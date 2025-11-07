package domain

import "github.com/google/uuid"

type EventType string

const (
	PlayerJoinedEvent      EventType = "JOINED_PLAYER"
	PlayerReconnectedEvent EventType = "PLAYER_RECONNECTED_EVENT"
	PlayerLeftEvent        EventType = "LEFT_PLAYER"
	GameStartedEvent       EventType = "STARTED_GAME"
	GameStateChangedEvent  EventType = "CHANGED_GAME_STATE"
	GameOverEvent          EventType = "GAME_OVER"
)

type Event struct {
	Type     EventType `json:"type"`
	CliendID uuid.UUID `json:"clientID"`
}
