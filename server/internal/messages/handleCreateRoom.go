package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/room"
	"server/internal/server"

	"github.com/google/uuid"
)

func (h *Handler) handleCreateRoom(client *server.Client, rawMsg json.RawMessage) error {
	var msg CreateRoomMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	if len(msg.Nickname) < 1 || len(msg.Nickname) > 20 {
		return SendError(client, ErrorMessage{
			RefType: CreateRoom,
			Field: "nickname",
			Reason: "Nickname should be 1 to 20 charecters long",
		})
	}

	newRoom := room.Room{
		ID:      uuid.New(),
		OwnerID: client.ID,
		Started: false,
	}
	client.RoomID = newRoom.ID

	if err := h.Cfg.Store.SetRoom(context.Background(), &newRoom); err != nil {
		return err
	}

	if err := h.Cfg.Presence.JoinRoom(context.Background(), newRoom.ID, client); err != nil {
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