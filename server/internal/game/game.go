package game

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/google/uuid"
)

func NewGame(cfg GameConfig) (*Game, error) {
	cards := make(map[uuid.UUID]Card)

	game := &Game{
		GameID:     uuid.New(),
		GameConfig: cfg,
		Cards:      &cards,
	}

	return game, nil
}

func (g *Game) GenerateCards() {
	cards := make([]Card, 0)
	combination := make([]string, 0)

	generateCombinations(g.GameConfig.Features, g.GameConfig.VariationsNumber, 0, combination, &cards)
	
	for i := range cards {
		cards[i].CardID = uuid.New()
		cards[i].IsVisible = false
		cards[i].IsDiscarded = false
		(*g.Cards)[cards[i].CardID] = cards[i]
		g.Deck = append(g.Deck, cards[i])
	}
}

func generateCombinations(features []Feature, variation int, index int, combination []string, deck *[]Card) {
	if index == len(features) {
		card := Card{}
		for i, feature := range features {
			value := combination[i]
			switch feature {
			case Color:
				card.Color = value
			case Shape:
				card.Shape = value
			case Number:
				card.Number = value
			case Shading:
				card.Shading = value
			case Rotation:
				card.Rotation = &value // optional field
			}
		}
		*deck = append(*deck, card)
		return
	}

	currentFeature := features[index]
	values := FeatureValues[currentFeature][:variation]

	for _, value := range values {
		combination = append(combination, value)
		generateCombinations(features, variation, index+1, combination, deck)
		combination = combination[:len(combination)-1]
	}
}

func (g *Game) ShuffleDeck() {
	rand.Shuffle(len(g.Deck), func(i, j int) {
		g.Deck[i], g.Deck[j] = g.Deck[j], g.Deck[i]
	})
}

func getFeatureValue(card Card, feature Feature) string {
	switch feature {
	case Color:
		return card.Color
	case Shape:
		return card.Shape
	case Number:
		return card.Number
	case Shading:
		return card.Shading
	case Rotation:
		if card.Rotation != nil {
			return *card.Rotation
		}
		return ""
	default:
		return ""
	}
}

func (g *Game) IsSet(cards []Card) bool {
	if len(cards) != g.GameConfig.VariationsNumber {
		return false
	}
	
	for _, feature := range g.GameConfig.Features {
	values := make(map[string]struct{})
	for _, card := range cards {
		val := getFeatureValue(card, feature)
		values[val] = struct{}{}
	}
	if len(values) != 1 && len(values) != g.GameConfig.VariationsNumber {
		return false
	}
}
	return true
}

func (g *Game) FindSet() []Card {
	c := make([]Card, 0)
	cardsInPlay := make([]Card, 0)
	for _, card := range g.Deck {
		if !card.IsVisible || card.IsDiscarded {
			continue
		}
		cardsInPlay = append(cardsInPlay, card)
	}

	var findCombinations func(startIndex int) []Card
	findCombinations = func(startIndex int) []Card {
		if (len(c) == g.GameConfig.VariationsNumber) {
			if g.IsSet(c) {
				return c
			}
			return nil
		}

		for i := startIndex; i < len(cardsInPlay); i++ {
			c = append(c, cardsInPlay[i])
			result := findCombinations(i + 1)
			if result != nil {
				return result
			}
			c = c[:len(c)-1]
		}
		return nil
	}

	return findCombinations(0)
}

func (g *Game) IsSetAvailable() bool {
	set := g.FindSet()
	return set != nil
}

func (g *Game) HandleCheckSet(ids []uuid.UUID) error {
	cards := make([]Card, 0)

	// validate cards exist
	for _, id := range ids {
		card, ok := (*g.Cards)[id]
		if !ok {
			return fmt.Errorf("card with given id doesn't exist")
		}
		cards = append(cards, card)
	}

	isSet := g.IsSet(cards)
	if !isSet {
		return nil
	}

	// update cards in map
	for i := range cards {
		cards[i].IsDiscarded = true
		cards[i].IsVisible = false
		(*g.Cards)[cards[i].CardID] = cards[i]
	}

	// update cards in deck
	counter := 0
	deck := g.Deck
	for i, card := range deck {
		if !slices.Contains(ids, card.CardID) {
			continue
		}
		deck[i].IsVisible = false
		deck[i].IsDiscarded = true
		counter++
		if counter == len(ids) {
			break
		}
	}
	g.Deck = deck

	return nil
}

func (g *Game) DealCards(n int) {
	dealt := 0
	for i, card := range g.Deck {
		if card.IsDiscarded || card.IsVisible {
			continue
		}
		card.IsVisible = true
		g.Deck[i] = card
		(*g.Cards)[card.CardID] = card
		dealt++

		if dealt == n {
			break
		}
	}
}