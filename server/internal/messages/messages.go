package messages

import (
	"encoding/json"
	"server/internal/config"
	"server/internal/game"
	"server/internal/server"

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
)

type StartGameMessage struct {
	InMessage
	GameVersion game.GameVersion `json:"gameVersion"`
	RoomID uuid.UUID
}

type JoinRoomMessage struct {
	InMessage
	RoomID uuid.UUID
}

type Handler struct {
	Cfg *config.Config
}

type InMessageHandler func(client *server.Client, rawMsg json.RawMessage) error

func (h *Handler) RegisterHandlers() map[InMessageType]InMessageHandler {
	return map[InMessageType]InMessageHandler{
		CreateRoom: h.handleCreateRoom,
		JoinRoom:   h.handleJoinRoom,
		StartGame:  h.handleStartGame,
	}
}

// Out

type OutMessageType string
type BaseOutMessage struct {
	Type OutMessageType `json:"type"`
}

const (
	CreatedRoom OutMessageType = "CREATED_ROOM"
	JoinedRoom  OutMessageType = "JOINED_ROOM"
	StartedGame OutMessageType = "STARTED_GAME"
)

type CreatedRoomMessage struct {
	BaseOutMessage
	RoomID uuid.UUID `json:"roomID"`
	PlayerID uuid.UUID `json:"playerID"`
}

type JoinedRoomMessage struct {
	BaseOutMessage
	RoomID uuid.UUID `json:"roomID"`
	PlayerID uuid.UUID `json:"playerID"`
	Error  string `json:"error,omitempty"`
}

type StartedGameMessage struct {
	BaseOutMessage
	GameID uuid.UUID   `json:"gameID"`
	Deck   []game.Card `json:"deck"`
}
