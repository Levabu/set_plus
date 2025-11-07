package presence

import (
	"context"

	"github.com/google/uuid"
)

type Presence interface {
	JoinRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error
	LeaveRoom(ctx context.Context, clientID uuid.UUID) error
	GetActiveRoomMembersIDs(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error)
	GetActiveRoomMembers(ctx context.Context, roomID uuid.UUID) ([]PresenceClient, error)
	SetClient(ctx context.Context, clientID uuid.UUID, status PresenceClient) error
	GetClient(ctx context.Context, clientID uuid.UUID) (PresenceClient, error)
	CleanupPresenceRoom(ctx context.Context, roomID uuid.UUID)
	RemoveClient(ctx context.Context, clientID uuid.UUID, roomID uuid.UUID) error
	// CleanupDisconnectedClients(ctx context.Context) error
	// IsClientConnected(ctx context.Context, clientID uuid.UUID) (bool, error)
	// UpdateHeartbeat(ctx context.Context, clientID uuid.UUID) error
	// BroadcastToRoom(ctx context.Context, roomID uuid.UUID, message interface{}, localClients domain.LocalClientManager) error
}

type PresenceClient struct {
	ID             uuid.UUID `json:"id"`
	RoomID         uuid.UUID `json:"roomID,omitempty"`
	Connected      bool      `json:"connected"`
	DisconnectedAt int64     `json:"lastSeen"` // Unix timestamp
}
