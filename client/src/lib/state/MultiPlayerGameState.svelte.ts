import { ROTATIONS, type GameVersion } from "$lib/engine/types";
import { type ChangedGameStateMessage, type CheckSetResultMessage, type GameOverMessage, type Player } from "$lib/ws/messages";
import { GameState } from "./GameState.svelte";

export class MultiPlayerGameState extends GameState {
  playerID: string = $state("");
  players: Record<string, Player> = $state({});
  score = $derived((() => {
    const player = this.players[this.playerID];
    return player ? player.score : 0;
  })())
  isOver: boolean = $state(false);
  winnerIDs: string[] = $derived.by(() => {
    if (!this.isOver) return [];
    let maxScore = -Infinity;
    for (const playerID in this.players) {
      const player = this.players[playerID];
      if (player.score > maxScore) {
        maxScore = player.score;
      }
    }
    return  Object.values(this.players).filter(player => player.score === maxScore).map(player => player.id);
  })

  constructor(gameVersion: GameVersion) {
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
  }
}