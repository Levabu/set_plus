import type { Card, Feature, GameVersion, VariationsNumber } from "$lib/engine/types";
import Game, { type GameOptions } from "$lib/engine/Game";

export class SinglePlayerGameState {
  features: Feature[];
  variationsNumber: VariationsNumber;
  game: Game;
  id: string = crypto.randomUUID();

  deck: Card[] = $state<Card[]>([]);
  selectedIds: Card["id"][] = $derived(this.deck.filter(card => card.isSelected).map(card => card.id))
  drawPile: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && !c.isVisible))
  inPlayCards: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && c.isVisible)) 
  isSetAvailable: boolean = $state(true);

  constructor(gameVersion: GameVersion) {
    this.features = gameVersion.features;
    this.variationsNumber = gameVersion.variationsNumber;
    this.game = new Game({ features: this.features, variationsNumber: this.variationsNumber } as GameOptions);
    this.deck = this.game.generateDeck();
    this.dealCards(gameVersion.initialDeal);

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
        this.dealCards(this.variationsNumber)
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
    const availableCards = this.deck.filter(c => !c.isDiscarded && !c.isVisible);
    
    let dealt = 0;
    for (let i = 0; i < availableCards.length && dealt < cardsNumber; i++) {
      availableCards[i].isVisible = true;
      dealt++;
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