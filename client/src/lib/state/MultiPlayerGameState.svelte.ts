import type Game from "$lib/engine/Game";
import { ROTATIONS, type GameVersion } from "$lib/engine/types";
import { OUT_MESSAGES, IN_MESSAGES, type StartGameMessage, type StartedGameMessage, type CreateRoomMessage, type CreatedRoomMessage, type JoinedRoomMessage } from "$lib/ws/messages";
import { CONNECTION_STATUS, WS } from "$lib/ws/ws.svelte";
import { GameState } from "./GameState.svelte";

export class MultiPlayerGameState extends GameState {
  // id = crypto.randomUUID();
  ws: WS;
  roomID: string = $state<string>("");
  playerID: string = $state<string>("");

  constructor(gameVersion: GameVersion) {
    super(gameVersion);
    this.ws = new WS("ws://localhost:8080/ws")
    
    $effect(() => {
      if (this.deck.length === 0 && this.ws.connectionStatus === CONNECTION_STATUS.CONNECTED) {
        this.ws.send({
          type: OUT_MESSAGES.CREATE_ROOM
        } as CreateRoomMessage);
      }
    })

    $effect(() => {
      if (this.ws.messages.length === 0) return;
      const lastMessage = this.ws.messages[this.ws.messages.length - 1];

      switch (lastMessage.type) {
        case IN_MESSAGES.CREATED_ROOM:
          // Handle room creation logic if needed
          console.log("Room created:", lastMessage);
          this.handleCreatedRoomMessage(lastMessage as CreatedRoomMessage);
          break;
        case IN_MESSAGES.JOINED_ROOM:
          // Handle joining room logic if needed
          console.log("Joined room:", lastMessage);
          this.handleJoinedRoomMessage(lastMessage as JoinedRoomMessage);
          break;
        case IN_MESSAGES.STARTED_GAME:
          this.handleGameCreatedMessage(lastMessage as StartedGameMessage);
          break;
        default:
          console.warn("Unhandled message type:");
          break;
      }

    })
  }

  handleCreatedRoomMessage(message: CreatedRoomMessage): void {
    console.log("Room created with ID:", message.roomID);
    this.roomID = message.roomID;
    this.playerID = message.playerID;

    // temp
    this.ws.send({
      type: OUT_MESSAGES.START_GAME,
      gameVersion: this.gameVersion.key,
      roomID: this.roomID
    } as StartGameMessage);
  }

  handleJoinedRoomMessage(message: JoinedRoomMessage): void {
    console.log("Joined room:", message.roomID, "as player:", message.playerID);
    this.roomID = message.roomID;
    this.playerID = message.playerID;

    if (message.error) {
      console.error("Error joining room:", message.error);
      // Handle error logic here, e.g., show a notification to the user
    }
  }

  handleGameCreatedMessage(message: StartedGameMessage): void {
    this.id = message.gameID;
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