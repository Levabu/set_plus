package events

import (
	"encoding/json"
	"log"
	"server/internal/domain"

	"github.com/google/uuid"
)

// This is called from Redis PubSub subscription
func (h *RoomEventHandler) HandleRoomEventMessage(roomID uuid.UUID, clientID uuid.UUID, msgData []byte) {
	var event domain.Event
	if err := json.Unmarshal(msgData, &event); err != nil {
		log.Printf("Failed to unmarshal room event: %v", err)
		return
	}

	if err := h.HandleRoomEvent(roomID, event); err != nil {
		log.Printf("Failed to handle room event: %v", err)
	}
}

func (h *RoomEventHandler) HandleRoomEvent(roomID uuid.UUID, event domain.Event) error {
	switch event.Type {
	case domain.PlayerJoinedEvent:
		return h.handleJoinedPlayer(roomID, event)
	case domain.GameStartedEvent:
		return h.handleStartedGame(roomID, event)
	case domain.GameStateChangedEvent:
		return h.handleChangedGameState(roomID, event)
	case domain.GameOverEvent:
		return h.handleGameOver(roomID, event)
	}
	return nil
}
