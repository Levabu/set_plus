package domain

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type MessageRouter interface {
	RegisterHandler(msgType InMessageType, handler MessageHandler)
	HandleMessage(client *Client, msgType InMessageType, rawMsg json.RawMessage) error
}

type EventPublisher interface {
	PublishRoomEvent(ctx context.Context, roomID uuid.UUID, event interface{}) error
	PublishGameEvent(ctx context.Context, gameID uuid.UUID, event interface{}) error
}

type ConnectionManager interface {
	HandleConnection(client *Client)
	BroadcastToRoom(roomID uuid.UUID, message interface{}) error
	SendToClient(clientID uuid.UUID, message interface{}) error
}