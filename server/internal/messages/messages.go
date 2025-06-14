package messages

import (
	"encoding/json"
	"server/internal/game"
	"server/internal/server"
)

type BaseMessage struct {
	Type string `json:"type"`
}

type CreateGameMessage struct {
	BaseMessage
	GameVersion game.GameVersion `json:"gameVersion"`
}

type MessageHandler func(client *server.Client, rawMsg json.RawMessage) error

var MessageHandlers = map[string]MessageHandler {
	"CREATE_GAME": handleCreateGame,
}
