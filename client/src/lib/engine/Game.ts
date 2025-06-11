import { featureValues, type Card, type Feature, type VariationsNumber } from "./types";

export interface GameOptions {
  features: Feature[];
  variationsNumber: VariationsNumber;
}

export class Game {
  readonly features: Feature[];
  readonly variationsNumber: VariationsNumber;

  constructor(options: GameOptions) {
    this.features = options.features;
    this.variationsNumber = options.variationsNumber;
  }

  private generateCombinations(
    index: number,
    combination: Card[],
    deck: Card[]
  ): void {
    if (index === this.features.length) {
      const card = combination.reduce(
        (acc, feature) => ({ ...acc, ...feature }),
        {} as Card
      );
      deck.push(card);
      return;
    }

    const currentFeature = this.features[index];
    const currentFeatureValues = featureValues[currentFeature].slice(
      0,
      this.variationsNumber
    );

    for (const value of currentFeatureValues) {
      combination.push({
        ...{ [currentFeature]: value },
        id: crypto.randomUUID(),
        isVisible: false,
        isSelected: false,
        isDiscarded: false,
      } as Card);
      this.generateCombinations(index + 1, combination, deck);
      combination.pop();
    }
  }

  generateDeck(): Card[] {
    const deck: Card[] = [];
    this.generateCombinations(0, [], deck);
    // Game.shuffleDeck(deck);
    for (let i = 0; i < deck.length; i++) {
      deck[i].id = crypto.randomUUID();
    }
    return deck;
  }

  static shuffleDeck(deck: Card[]): void {
    for (let i = deck.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [deck[i], deck[j]] = [deck[j], deck[i]];
    }
  }

  static isSet(cards: Card[], features: Feature[], variationsNumber: VariationsNumber): boolean {
    if (cards.length !== variationsNumber) return false;

    for (const feature of features) {
      const values = new Set(cards.map(card => card[feature]));
      if (values.size !== 1 && values.size !== variationsNumber) return false;
    }
    return true;
  }
  
  static findSet(cards: Card[], features: Feature[], variationsNumber: VariationsNumber): Card[] | null {
    function findCombinations(
      startIndex: number,
      currentCombination: Card[]
    ): Card[] | null {
      if (currentCombination.length === variationsNumber) {
        if (Game.isSet(currentCombination, features, variationsNumber)) {
          return currentCombination
        }
        return null
      }
      
      for (let i = startIndex; i < cards.length; i++) {
        currentCombination.push(cards[i]);
        const result = findCombinations(i + 1, currentCombination);
        if (result) {
          return result;
        }
        currentCombination.pop();
      }
      return null;
    }
    return findCombinations(0, [])
  }
  
  static isSetAvailable(cards: Card[], features: Feature[], variationsNumber: VariationsNumber): boolean {
    const set = this.findSet(cards, features, variationsNumber)
    return set === null ? false : true
  }
}

export default Game;
