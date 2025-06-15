package room

import (
	// "context"
	// "encoding/json"
	// "log"
	// "server/internal/messages"

	"github.com/google/uuid"
)

type Room struct {
	ID        uuid.UUID
	OwnerID   uuid.UUID
	GameID	  uuid.UUID
	Started   bool
}

// func HandleRoomEvent(handler *messages.Handler, id uuid.UUID, rawMsg []byte) {
// 	var event Event
// 	if err := json.Unmarshal(rawMsg, &event); err != nil {
// 		log.Println("invalid event payload:", err)
// 		return
// 	}

// 	switch event.Type {
// 	case PlayerJoined:

// 	}
// }