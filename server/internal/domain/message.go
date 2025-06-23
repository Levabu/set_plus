package domain

import (
	"encoding/json"
	"server/internal/game"

	"github.com/google/uuid"
)

type InMessageType string
type InMessage struct {
	Type InMessageType `json:"type"`
}

const (
	CreateRoom InMessageType = "CREATE_ROOM"
	JoinRoom   InMessageType = "JOIN_ROOM"
	StartGame  InMessageType = "START_GAME"
	CheckSet   InMessageType = "CHECK_SET"
)

type StartGameMessage struct {
	InMessage
	GameVersion game.GameVersion `json:"gameVersion"`
	RoomID      uuid.UUID        `json:"roomID"`
}

type CreateRoomMessage struct {
	InMessage
	Nickname string `json:"nickname"`
}

type JoinRoomMessage struct {
	InMessage
	RoomID   uuid.UUID `json:"roomID"`
	Nickname string    `json:"nickname"`
}

type CheckSetMessage struct {
	InMessage
	CardIDs  []uuid.UUID `json:"cardIDs"`
	PlayerID uuid.UUID   `json:"playerID"`
	RoomID   uuid.UUID   `json:"roomID"`
	GameID   uuid.UUID   `json:"gameID"`
}

type OutMessageType string
type BaseOutMessage struct {
	Type OutMessageType `json:"type"`
}

const (
	CreatedRoom      OutMessageType = "CREATED_ROOM"
	JoinedRoom       OutMessageType = "JOINED_ROOM"
	StartedGame      OutMessageType = "STARTED_GAME"
	CheckSetResult   OutMessageType = "CHECK_SET_RESULT"
	ChangedGameState OutMessageType = "CHANGED_GAME_STATE"
	GameOver         OutMessageType = "GAME_OVER"
	ErrorOut         OutMessageType = "ERROR"
)

type CreatedRoomMessage struct {
	BaseOutMessage
	RoomID   uuid.UUID `json:"roomID"`
	PlayerID uuid.UUID `json:"playerID"`
	Nickname string    `json:"nickname"`
}

type JoinedRoomMessage struct {
	BaseOutMessage
	RoomID   uuid.UUID `json:"roomID"`
	PlayerID uuid.UUID `json:"playerID"`
	Nickname string    `json:"nickname"`
}

type StartedGameMessage struct {
	BaseOutMessage
	GameID      uuid.UUID                 `json:"gameID"`
	GameVersion game.GameVersion          `json:"gameVersion"`
	Deck        []game.Card               `json:"deck"`
	Players     map[uuid.UUID]game.Player `json:"players"`
}

type CheckSetResultMessage struct {
	BaseOutMessage
	IsSet bool `json:"isSet"`
}

type ChangedGameStateMessage struct {
	BaseOutMessage
	GameID  uuid.UUID                 `json:"gameID"`
	Deck    []game.Card               `json:"deck"`
	Players map[uuid.UUID]game.Player `json:"players"`
}

type GameOverMessage struct {
	BaseOutMessage
	GameID  uuid.UUID                 `json:"gameID"`
	Deck    []game.Card               `json:"deck"`
	Players map[uuid.UUID]game.Player `json:"players"`
}

type ErrorMessage struct {
	RefType InMessageType `json:"refType"`
	Field   string        `json:"field"`
	Reason  string        `json:"reason"`
}

type MessageHandler func(client *Client, rawMsg json.RawMessage) error