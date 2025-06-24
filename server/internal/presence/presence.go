package presence

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

type Presence interface {
	JoinRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error
	LeaveRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error
	GetRoomMembers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error)
	SubscribeToRoom(ctx context.Context, roomID uuid.UUID, handler func(clientID uuid.UUID, msg []byte)) error
	SetClientStatus(ctx context.Context, clientID uuid.UUID, status ClientStatus) error
	GetClientStatus(ctx context.Context, clientID uuid.UUID) (ClientStatus, error)
	CleanupDisconnectedClients(ctx context.Context) error
	IsClientConnected(ctx context.Context, clientID uuid.UUID) (bool, error)
	UpdateHeartbeat(ctx context.Context, clientID uuid.UUID) error
	BroadcastToRoom(ctx context.Context, roomID uuid.UUID, message interface{}, localClients domain.ClientManager) error
}

type ClientStatus struct {
	ID          uuid.UUID `json:"id"`
	RoomID      uuid.UUID `json:"roomID,omitempty"`
	// Nickname    string    `json:"nickname"`
	Connected   bool      `json:"connected"`
	LastSeen    int64     `json:"lastSeen"` // Unix timestamp
	Reconnected bool      `json:"reconnected"`
}