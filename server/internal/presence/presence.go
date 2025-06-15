package presence

import (
	"context"

	"github.com/google/uuid"
)

type Presence interface {
	JoinRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error
	LeaveRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error
	GetRoomMembers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error)
	// BroadcastToRoom(ctx context.Context, roomID uuid.UUID, message []byte) error
	SubscribeToRoom(ctx context.Context, roomID uuid.UUID, handler func(clientID uuid.UUID, msg []byte)) error
}