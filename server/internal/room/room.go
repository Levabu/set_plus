package room

import (
	// "context"
	// "encoding/json"
	// "log"
	// "server/internal/messages"

	"github.com/google/uuid"
)

type Room struct {
	ID       uuid.UUID
	OwnerID  uuid.UUID
	GameID   uuid.UUID
	Started  bool
}
