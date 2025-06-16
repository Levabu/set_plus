package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/room"
	"server/internal/server"
)

func (h *Handler) handleJoinRoom(client *server.Client, rawMsg json.RawMessage) error {
	var msg JoinRoomMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	joinedRoom, err := h.Cfg.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return err
	}

	if joinedRoom.Started {
		return SendError(client, ErrorMessage{
			RefType: JoinRoom,
			Reason: "game already started",
		})
	}

	if err := h.Cfg.Presence.JoinRoom(context.Background(), joinedRoom.ID, client); err != nil {
		return err
	}

	err = h.Cfg.Store.PublishRoomUpdate(context.Background(), joinedRoom.ID, room.Event{
		Type:     room.JoinedPlayer,
		CliendID: client.ID,
	})
	if err != nil {
		return err
	}

	return nil
}