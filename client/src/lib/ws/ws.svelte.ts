import { GameVersions, ROTATIONS, type GameVersion, type GameVersionKey } from "$lib/engine/types";
import { MultiPlayerGameState } from "$lib/state/MultiPlayerGameState.svelte";
import { type OutMessage, type InMessage, OUT_MESSAGES, type StartGameMessage, IN_MESSAGES, type CreatedRoomMessage, type JoinedRoomMessage, type StartedGameMessage, type CheckSetResultMessage, type ChangedGameStateMessage, type GameOverMessage, type CheckSetMessage } from "./messages";

export const CONNECTION_STATUS = {
  CONNECTED: 'connected',
  DISCONNECTED: 'disconnected',
  ERROR: 'error',
} as const;

export class WS {
  socket: WebSocket | null = null;
  messages: InMessage[] = $state([]);
  connectionStatus: string = $state(CONNECTION_STATUS.DISCONNECTED);
  game = $state<MultiPlayerGameState | null>(null);
  roomID: string = $state<string>("");
  playerID: string = $state<string>("");
  isRoomOwner: boolean = $state(false);

  constructor(url: string = "ws://localhost:8080") {
    let ws = new WebSocket(url)
    this.socket = ws

    ws.onopen = () => {
        this.socket = ws;
        this.connectionStatus = CONNECTION_STATUS.CONNECTED;
        console.log('WebSocket connected');
    };
    
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        data.isProcessed = false
        console.log('WebSocket message received:', data);
        this.messages.push(data);
    };
    
    ws.onclose = () => {
        this.socket = null;
        this.connectionStatus = CONNECTION_STATUS.DISCONNECTED;
        console.log('WebSocket disconnected');
    };
    
    ws.onerror = (error) => {
        this.connectionStatus = CONNECTION_STATUS.ERROR;
        console.error('WebSocket error:', error);
    };

    $effect(() => {
      if (this.messages.length === 0) return
      const lastMessage = this.messages[this.messages.length - 1]
      if (lastMessage.isProcessed) {
        return;
      }

      switch (lastMessage.type) {
        case IN_MESSAGES.CREATED_ROOM:
          this.handleCreatedRoomMessage(lastMessage as CreatedRoomMessage);
          break;
        case IN_MESSAGES.JOINED_ROOM:
          this.handleJoinedRoomMessage(lastMessage as JoinedRoomMessage);
          break;
        case IN_MESSAGES.STARTED_GAME:
          this.handleStartedGameMessage(lastMessage as StartedGameMessage);
          break;
        case IN_MESSAGES.CHECK_SET_RESULT:
          this.game?.handleCheckSetResultMessage(lastMessage as CheckSetResultMessage);
          break;
        case IN_MESSAGES.CHANGED_GAME_STATE:
          this.game?.handleGameStateUpdate(lastMessage as ChangedGameStateMessage);
          break;
        case IN_MESSAGES.GAME_OVER:
          this.game?.handleGameOverMessage(lastMessage as GameOverMessage);
          break;
        default:
          console.warn("Unhandled message type:");
          break;
      }
      lastMessage.isProcessed = true;
    })

    $effect(() => {
      if (!this.game) return
      if (this.game.selectedIds.length !== this.game.variationsNumber) return;
      this.send({
        type: OUT_MESSAGES.CHECK_SET,
        cardIDs: this.game.selectedIds,
        playerID: this.game.playerID,
        roomID: this.roomID,
        gameID: this.game.id
      } as CheckSetMessage);
    });
  }

  handleCreatedRoomMessage(message: CreatedRoomMessage): void {
    console.log("Room created with ID:", message.roomID);
    this.roomID = message.roomID;
    this.playerID = message.playerID;
    this.isRoomOwner = true
  }

  handleJoinedRoomMessage(message: JoinedRoomMessage): void {
    console.log("Joined room:", message.roomID, "as player:", message.playerID);
    if (this.playerID) return
    this.roomID = message.roomID;
    this.playerID = message.playerID;

    if (message.error) {
      console.error("Error joining room:", message.error);
    }
  }

  handleStartedGameMessage(message: StartedGameMessage): void {
    if (this.game) return
    const gameVersion = GameVersions[message.gameVersion]
    if (!gameVersion) return
    this.game = new MultiPlayerGameState(gameVersion)
    this.game.id = message.gameID;
    this.game.hasGameStarted = true;
    this.game.playerID = this.playerID

    this.game.deck = message.deck.map((card: any) => ({
      id: card.id,
      isVisible: card.isVisible,
      isSelected: card.isSelected,
      isDiscarded: card.isDiscarded,
      color: card.color,
      shape: card.shape,
      number: Number(card.number),
      shading: card.shading,
      rotation: card.rotation || ROTATIONS.vertical,
    }));
    this.game.players = message.players
  }

  handleStartGame(gameVersion: GameVersionKey): void {
    if (this.game) return
    if (this.isRoomOwner) {
      this.send({
        type: OUT_MESSAGES.START_GAME,
        gameVersion: gameVersion,
        roomID: this.roomID
      } as StartGameMessage);
    }
  }

  send(message: OutMessage) {
    console.log("Sending message:", message);
    if (this.socket && this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(JSON.stringify(message));
    } else {
      console.error("WebSocket is not open. Cannot send message:", this.socket?.readyState, message);
    }
  }
}