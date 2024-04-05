import { Card, Feature, OptionsNumber, featureValues } from "./features";

interface GameOptions {
  features: Feature[];
  optionsNumber: OptionsNumber;
}

class Game {
  readonly features: Feature[];
  readonly optionsNumber: OptionsNumber;
  deck: Card[];

  constructor(options: GameOptions) {
    this.features = options.features;
    this.optionsNumber = options.optionsNumber;
    this.deck = this.generateDeck();
  }

  generateCombinations(index: number, combination: Card[]): void {
    if (index === this.features.length) {
      const card = combination.reduce((acc, feature) => ({ ...acc, ...feature }), {});
      this.deck.push(card);
      return;
    }

    const currentFeature = this.features[index];
    const currentFeatureValues = featureValues[currentFeature].slice(0, this.optionsNumber);

    for (const value of currentFeatureValues) {
      combination.push({ [currentFeature]: value });
      this.generateCombinations(index + 1, combination);
      combination.pop();
    }
  }

  generateDeck(): Card[] {
    this.deck = [];
    this.generateCombinations(0, []);
    return this.deck;
  }
}

export default Game;
