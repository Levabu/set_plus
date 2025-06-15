package messages

import (
	"context"
	"encoding/json"
	"server/internal/room"
	"server/internal/server"

	"github.com/google/uuid"
)

func (h *Handler) handleCreateRoom(client *server.Client, rawMsg json.RawMessage) error {
	newRoom := room.Room{
		ID:      uuid.New(),
		OwnerID: client.ID,
		Started: false,
	}
	client.RoomID = newRoom.ID

	if err := h.Cfg.Store.SetRoom(context.Background(), &newRoom); err != nil {
		return err
	}

	if err := h.Cfg.Presence.JoinRoom(context.Background(), newRoom.ID, client.ID); err != nil {
		return err
	}

	go h.Cfg.Presence.SubscribeToRoom(context.Background(), newRoom.ID, func(clientID uuid.UUID, event room.Event) {
		h.HandleRoomEvent(clientID, event)
	})

	SendJSON(client, CreatedRoomMessage{
		BaseOutMessage: BaseOutMessage{Type: CreatedRoom},
		RoomID:         newRoom.ID,
		PlayerID:       newRoom.OwnerID,
	})
	return nil
}