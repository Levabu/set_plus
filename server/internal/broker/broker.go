package broker

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Broker interface {
	PublishRoomUpdate(ctx context.Context, roomID uuid.UUID, event domain.Event) error
	SubscribeToRoom(ctx context.Context, roomID uuid.UUID, handler func(clientID uuid.UUID, msg any)) error
}