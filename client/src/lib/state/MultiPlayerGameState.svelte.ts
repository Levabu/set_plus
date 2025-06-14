import type Game from "$lib/engine/Game";
import type { GameVersion } from "$lib/engine/types";
import { CLIENT_TO_SERVER_MESSAGES, type CreateGameMessage, type GameCreatedMessage } from "$lib/ws/messages";
import { CONNECTION_STATUS, WS } from "$lib/ws/ws.svelte";
import { GameState } from "./GameState.svelte";

export class MultiPlayerGameState extends GameState {
  // id = crypto.randomUUID();
  ws: WS;

  constructor(gameVersion: GameVersion) {
    super(gameVersion);
    this.ws = new WS("ws://localhost:8080/ws")
    
    $effect(() => {
      if (this.deck.length === 0 && this.ws.connectionStatus === CONNECTION_STATUS.CONNECTED) {
        this.ws.send({
          type: CLIENT_TO_SERVER_MESSAGES.CREATE_GAME,
          gameVersion: gameVersion.key
        } as CreateGameMessage);
      }
    })

    $effect(() => {
      if (this.ws.messages.length === 0) return;
      const lastMessage = this.ws.messages[this.ws.messages.length - 1];


    })
  }

  handleGameCreatedMessage(message: GameCreatedMessage): void {
    this.id = message.gameID;
    // Convert message.deck to Card[] if necessary
    this.deck = message.deck.map((card: any) => ({
      id: card.id,
      isVisible: card.isVisible,
      isSelected: card.isSelected,
      isDiscarded: card.isDiscarded
    }));
    this.dealCards(this.game.variationsNumber);
  }


}