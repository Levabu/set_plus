import { ROTATIONS, type GameVersion } from "$lib/engine/types";
import { OUT_MESSAGES, IN_MESSAGES, type StartGameMessage, type StartedGameMessage, type CreateRoomMessage, type CreatedRoomMessage, type JoinedRoomMessage } from "$lib/ws/messages";
import { CONNECTION_STATUS, WS } from "$lib/ws/ws.svelte";
import { GameState } from "./GameState.svelte";

export class MultiPlayerGameState extends GameState {
  ws: WS;
  roomID: string = $state<string>("");
  playerID: string = $state<string>("");
  private hasGameStarted: boolean = $state<boolean>(false);

  constructor(gameVersion: GameVersion) {
    console.log("MultiPlayerGameState constructor called with gameVersion:", gameVersion);
    super(gameVersion);
    this.ws = new WS("ws://localhost:8080/ws")
    
    $effect(() => {
      if (this.roomID === "" && this.ws.connectionStatus === CONNECTION_STATUS.CONNECTED) {
        this.ws.send({
          type: OUT_MESSAGES.CREATE_ROOM
        } as CreateRoomMessage);
      }
    })

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
          this.handleGameCreatedMessage(lastMessage as StartedGameMessage);
          break;
        default:
          console.warn("Unhandled message type:");
          break;
      }
      lastMessage.isProcessed = true;
    })
  }

  handleCreatedRoomMessage(message: CreatedRoomMessage): void {
    console.log("Room created with ID:", message.roomID);
    this.roomID = message.roomID;
    this.playerID = message.playerID;

    // temp
    if (!this.hasGameStarted) {
      this.ws.send({
        type: OUT_MESSAGES.START_GAME,
        gameVersion: this.gameVersion.key,
        roomID: this.roomID
      } as StartGameMessage);
    }
  }

  handleJoinedRoomMessage(message: JoinedRoomMessage): void {
    console.log("Joined room:", message.roomID, "as player:", message.playerID);
    this.roomID = message.roomID;
    this.playerID = message.playerID;

    if (message.error) {
      console.error("Error joining room:", message.error);
    }
  }

  handleGameCreatedMessage(message: StartedGameMessage): void {
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
  }


}