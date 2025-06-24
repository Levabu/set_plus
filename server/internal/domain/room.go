package domain

import (
	"github.com/google/uuid"
)

type Room struct {
	ID       uuid.UUID
	OwnerID  uuid.UUID
	GameID   uuid.UUID
	Started  bool
}
