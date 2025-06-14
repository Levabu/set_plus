import type { Card, Feature, GameVersion, VariationsNumber } from "$lib/engine/types";
import Game, { type GameOptions } from "$lib/engine/Game";

export class GameState {
  features: Feature[];
  variationsNumber: VariationsNumber;
  game: Game;
  id!: string;

  deck: Card[] = $state<Card[]>([]);
  selectedIds: Card["id"][] = $derived(this.deck.filter(card => card.isSelected).map(card => card.id))
  drawPile: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && !c.isVisible))
  inPlayCards: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && c.isVisible)) 

  constructor(gameVersion: GameVersion) {
    this.features = gameVersion.features;
    this.variationsNumber = gameVersion.variationsNumber;
    this.game = new Game({ features: this.features, variationsNumber: this.variationsNumber } as GameOptions);
  }

  toggleSelectCard(cardId: Card["id"]): void {
    const card = this.deck.find(c => c.id === cardId);
    if (card) {
      card.isSelected = !card.isSelected;
    }
  }

  resetSelectedCards(): void {
    for (const card of this.deck) {
      card.isSelected = false;
    }
  }

  findSet(): Card[] | null {
    return Game.findSet(
      this.inPlayCards,
      this.features,
      this.variationsNumber
    );
  }
}