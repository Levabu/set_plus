package events

import (
	"context"
	"encoding/json"
	"log"
	"server/internal/config"
	"server/internal/domain"

	"github.com/google/uuid"
)

type RoomEventHandler struct {
	config *config.Config
}

func NewRoomEventHandler(config *config.Config) *RoomEventHandler {
	return &RoomEventHandler{
		config: config,
	}
}

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
	case domain.PlayerLeftEvent:
		return h.handleDisconnectedPlayer(roomID, event)
	case domain.PlayerReconnectedEvent:
		return h.handleReconnectedPlayer(roomID, event)
	case domain.GameStartedEvent:
		return h.handleStartedGame(roomID, event)
	case domain.GameStateChangedEvent:
		return h.handleChangedGameState(roomID, event)
	case domain.GameOverEvent:
		return h.handleGameOver(roomID, event)
	}
	return nil
}

func (h *RoomEventHandler) BroadcastToRoom(ctx context.Context, roomID uuid.UUID, message interface{}, localClients domain.LocalClientManager) error {
	members, err := h.config.Presence.GetActiveRoomMembersIDs(ctx, roomID)
	if err != nil {
		return err
	}

	for _, memberID := range members {
		memberClient := localClients.Get(memberID)
		if memberClient != nil {
			domain.SendJSON(memberClient, message)
		}
	}
	return nil
}
