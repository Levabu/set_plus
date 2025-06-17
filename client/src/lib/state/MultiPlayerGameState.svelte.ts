import { ROTATIONS, type GameVersion } from "$lib/engine/types";
import { OUT_MESSAGES, IN_MESSAGES, type StartGameMessage, type StartedGameMessage, type CreateRoomMessage, type CreatedRoomMessage, type JoinedRoomMessage, type CheckSetResultMessage, type CheckSetMessage, type ChangedGameStateMessage, type Player, type GameOverMessage } from "$lib/ws/messages";
import { CONNECTION_STATUS, WS } from "$lib/ws/ws.svelte";
import { GameState } from "./GameState.svelte";

export class MultiPlayerGameState extends GameState {
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
  }

  handleCheckSetResultMessage(message: CheckSetResultMessage): void {
    const isSet = message.isSet;
    if (!isSet) {
      return
    }
    this.resetSelectedCards();
  }

  handleGameStateUpdate(message: ChangedGameStateMessage): void {
    if (message.gameID !== this.id) {
      console.warn("Received game state update for a different game ID:", message.gameID);
      return;
    }
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
    console.log("Game over! Final scores:", message.players);
    alert("Game Over!")
  }
}