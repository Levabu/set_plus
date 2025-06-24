package events

import (
	"context"
	"server/internal/config"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Publisher struct {
	connectionManager domain.ConnectionManager
}

func NewPublisher(connectionManager domain.ConnectionManager) *Publisher {
	return &Publisher{
		connectionManager: connectionManager,
	}
}

func (p *Publisher) PublishRoomEvent(ctx context.Context, roomID uuid.UUID, event interface{}) error {
	return p.connectionManager.BroadcastToRoom(roomID, event)
}

func (p *Publisher) PublishGameEvent(ctx context.Context, gameID uuid.UUID, event interface{}) error {
	return nil
}

type RoomEventHandler struct {
	config *config.Config
}

func NewRoomEventHandler(config *config.Config) *RoomEventHandler {
	return &RoomEventHandler{
		config: config,
	}
}


