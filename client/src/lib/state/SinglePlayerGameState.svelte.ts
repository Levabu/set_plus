import type { Card, Feature, VariationsNumber } from "$lib/engine/types";
import Game, { type GameOptions } from "$lib/engine/Game";

export class SinglePlayerGameState {
  readonly features: Feature[];
  readonly variationsNumber: VariationsNumber;
  readonly game: Game;

  deck: Card[] = $state<Card[]>([]);
  selectedIds: Card["id"][] = $derived(this.deck.filter(card => card.isSelected).map(card => card.id))
  drawPile: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && !c.isVisible))
  inPlayCards: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && c.isVisible)) 
  isSetAvailable: boolean = $state(true);

  constructor(features: Feature[], variationsNumber: VariationsNumber) {
    this.features = features;
    this.variationsNumber = variationsNumber;
    this.game = new Game({ features, variationsNumber } as GameOptions);
    this.deck = this.game.generateDeck();
    this.dealCards(this.variationsNumber * 4);

    $effect(() => {
      if (this.selectedIds.length !== this.variationsNumber) return;
      const isSet = this.isSelectedSet();
      if (isSet) {
        this.removeSelectedCards();
        this.dealCards(this.variationsNumber)
      } else {
        this.resetSelectedCards();
      }
  });

    $effect(() => {
      this.isSetAvailable = Game.isSetAvailable(
        this.inPlayCards,
        this.features,
        this.variationsNumber
      )
    })

    $effect(() => {
      if (!this.isSetAvailable) {
        this.dealCards(1)
      }
    })
  }

  toggleSelectCard(cardId: Card["id"]): void {
    const card = this.deck.find(c => c.id === cardId);
    if (card) {
      card.isSelected = !card.isSelected;
    }
  }

  isSelectedSet(): boolean {
    return Game.isSet(
      this.deck.filter(card => card.isSelected),
      this.game.features,
      this.game.variationsNumber
    );
  }

  removeSelectedCards(): void {
    for (const card of this.deck) {
      if (card.isSelected) {
        card.isVisible = false;
        card.isSelected = false;
        card.isDiscarded = true;
      }
    }
  }

  resetSelectedCards(): void {
    for (const card of this.deck) {
      card.isSelected = false;
    }
  }
  
  dealCards(cardsNumber: number): void {
    let dealt = 0;
    while (dealt < cardsNumber && dealt < this.drawPile.length) {
      this.drawPile[0].isVisible = true
      dealt++
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