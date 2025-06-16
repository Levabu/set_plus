import { ROTATIONS, type GameVersion } from "$lib/engine/types";
import { OUT_MESSAGES, IN_MESSAGES, type StartGameMessage, type StartedGameMessage, type CreateRoomMessage, type CreatedRoomMessage, type JoinedRoomMessage, type CheckSetResultMessage, type CheckSetMessage, type ChangedGameStateMessage, type Player, type GameOverMessage } from "$lib/ws/messages";
import { CONNECTION_STATUS, WS } from "$lib/ws/ws.svelte";
import { GameState } from "./GameState.svelte";

export class MultiPlayerGameState extends GameState {
  ws: WS;
  roomID: string = $state<string>("");
  isRoomOwner: boolean = $state(false);
  playerID: string = $state<string>("");
  players: Record<string, Player> = $state<Record<string, Player>>({});
  hasGameStarted: boolean = $state<boolean>(false);
  score = $derived((() => {
    const player = this.players[this.playerID];
    return player ? player.score : 0;
  })())
  isOver: boolean = $state<boolean>(false);

  constructor(gameVersion: GameVersion) {
    console.log("MultiPlayerGameState constructor called with gameVersion:", gameVersion);
    super(gameVersion);
    this.ws = new WS("ws://localhost:8080/ws")
    
    // $effect(() => {
    //   if (this.roomID === "" && this.ws.connectionStatus === CONNECTION_STATUS.CONNECTED) {
    //     this.ws.send({
    //       type: OUT_MESSAGES.CREATE_ROOM
    //     } as CreateRoomMessage);
    //   }
    // })

    $effect(() => {
      if (this.ws.messages.length === 0) return
      const lastMessage = this.ws.messages[this.ws.messages.length - 1]
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
          this.handleCheckSetResultMessage(lastMessage as CheckSetResultMessage);
          break;
        case IN_MESSAGES.CHANGED_GAME_STATE:
          this.handleGameStateUpdate(lastMessage as ChangedGameStateMessage);
          break;
        case IN_MESSAGES.GAME_OVER:
          this.handleGameOverMessage(lastMessage as GameOverMessage);
          break;
        default:
          console.warn("Unhandled message type:");
          break;
      }
      lastMessage.isProcessed = true;
    })

    $effect(() => {
      if (this.selectedIds.length !== this.variationsNumber) return;
      this.ws.send({
        type: OUT_MESSAGES.CHECK_SET,
        cardIDs: this.selectedIds,
        playerID: this.playerID,
        roomID: this.roomID,
        gameID: this.id
      } as CheckSetMessage);
    });
  }

  handleCreatedRoomMessage(message: CreatedRoomMessage): void {
    console.log("Room created with ID:", message.roomID);
    this.roomID = message.roomID;
    this.playerID = message.playerID;
    this.isRoomOwner = true

    // temp
    // if (!this.hasGameStarted) {
    //   this.ws.send({
    //     type: OUT_MESSAGES.START_GAME,
    //     gameVersion: this.gameVersion.key,
    //     roomID: this.roomID
    //   } as StartGameMessage);
    // }
  }

  handleStartGame(): void {
    if (!this.hasGameStarted && this.isRoomOwner) {
      this.ws.send({
        type: OUT_MESSAGES.START_GAME,
        gameVersion: this.gameVersion.key,
        roomID: this.roomID
      } as StartGameMessage);
    }
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
    this.id = message.gameID;
    this.hasGameStarted = true;
    // Convert message.deck to Card[] if necessary
    this.deck = message.deck.map((card: any) => ({
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
    this.players = message.players
  }

  handleCheckSetResultMessage(message: CheckSetResultMessage): void {
    const isSet = message.isSet;
    if (!isSet) {
      return
    }
    // this.score += 1;
    this.resetSelectedCards();
  }

  handleGameStateUpdate(message: ChangedGameStateMessage): void {
    if (message.gameID !== this.id) {
      console.warn("Received game state update for a different game ID:", message.gameID);
      return;
    }
    // Update the game state based on the message
    this.deck = message.deck.map((card: any) => ({
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
    this.players = message.players
  }

  handleGameOverMessage(message: GameOverMessage): void {
    this.isOver = true;
    // Handle game over logic, e.g., display final scores
    console.log("Game over! Final scores:", message.players);
  }
}