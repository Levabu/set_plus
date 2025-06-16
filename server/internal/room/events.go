package room

import "github.com/google/uuid"

type EventType string

const (
	JoinedPlayer     EventType = "JOINED_PLAYER"
	StartedGame      EventType = "STARTED_GAME"
	ChangedGameState EventType = "CHANGED_GAME_STATE"
	GameOver         EventType = "GAME_OVER"
)

type Event struct {
	Type     EventType `json:"type"`
	CliendID uuid.UUID `json:"clientID"`
}
