import type { Card, Feature, GameVersion, VariationsNumber } from "$lib/engine/types";
import Game, { type GameOptions } from "$lib/engine/Game";
import { GameState } from "./GameState.svelte";

export class SinglePlayerGameState extends GameState {
  id: string = crypto.randomUUID();

  // deck: Card[] = $state<Card[]>([]);
  // selectedIds: Card["id"][] = $derived(this.deck.filter(card => card.isSelected).map(card => card.id))
  // drawPile: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && !c.isVisible))
  // inPlayCards: Card[] = $derived(this.deck.filter(c => !c.isDiscarded && c.isVisible)) 
  isSetAvailable: boolean = $state(true);

  constructor(gameVersion: GameVersion) {
    super(gameVersion);
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
  
  dealCards(cardsNumber: number): void {
    const availableCards = this.deck.filter(c => !c.isDiscarded && !c.isVisible);
    
    let dealt = 0;
    for (let i = 0; i < availableCards.length && dealt < cardsNumber; i++) {
      availableCards[i].isVisible = true;
      dealt++;
    }
  }
}