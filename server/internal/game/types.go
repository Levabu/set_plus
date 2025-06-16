package game

import "github.com/google/uuid"

type Feature string

const (
	Color    Feature = "color"
	Shape    Feature = "shape"
	Number   Feature = "number"
	Shading  Feature = "shading"
	Rotation Feature = "rotation"
)

var FeatureValues = map[Feature][]string{
	Color:    {"c1", "c2", "c3", "c4"},
	Shape:    {"diamond", "squiggle", "oval", "arrow"},
	Number:   {"1", "2", "3", "4"},
	Shading:  {"solid", "striped", "empty", "dotted"},
	Rotation: {"vertical", "horizontal", "diagonal"},
}

// type Card struct {
// 	CardID uuid.UUID `json:"id"`
// 	Features map[Feature]string `json:"features"`
// 	IsVisible bool `json:"isVisible"`
// 	IsDiscarded bool `json:"isDiscarded"`
// }

type Card struct {
	CardID      uuid.UUID `json:"id"`
	Color       string    `json:"color"`
	Shape       string    `json:"shape"`
	Number      string    `json:"number"`
	Shading     string    `json:"shading"`
	Rotation    *string   `json:"rotation,omitempty"`
	IsVisible   bool      `json:"isVisible"`
	IsDiscarded bool      `json:"isDiscarded"`
}

type GameVersion string

const (
	Classic GameVersion = "classic"
	V5x3    GameVersion = "v5x3"
	V4x4    GameVersion = "v4x4"
)

func (v GameVersion) IsValid() bool {
	switch v {
	case Classic, V4x4, V5x3:
		return true
	default:
		return false
	}
}

type GameConfig struct {
	Features         []Feature
	VariationsNumber int
	InitialDeal      int
}

var GameVersions = map[GameVersion]GameConfig{
	Classic: {
		Features:         []Feature{Color, Shape, Number, Shading},
		VariationsNumber: 3,
		InitialDeal:      12,
	},
	V5x3: {
		Features:         []Feature{Color, Shape, Number, Shading, Rotation},
		VariationsNumber: 3,
		InitialDeal:      15,
	},
	V4x4: {
		Features:         []Feature{Color, Shape, Number, Shading},
		VariationsNumber: 4,
		InitialDeal:      16,
	},
}

type Player struct {
	ID       uuid.UUID `json:"id"`
	Nickname string    `json:"nickname"`
	Score    int       `json:"score"`
}

type Game struct {
	GameID      uuid.UUID
	GameVersion GameVersion
	GameConfig  GameConfig
	Cards       *map[uuid.UUID]Card
	Deck        []Card
	Players     *map[uuid.UUID]Player
	Finished    bool
}
