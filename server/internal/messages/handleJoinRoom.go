package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/room"
	"server/internal/server"
)

func (h *Handler) handleJoinRoom(client *server.Client, rawMsg json.RawMessage) error {
	var msg JoinRoomMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	if len(msg.Nickname) < 1 || len(msg.Nickname) > 20 {
		return SendError(client, ErrorMessage{
			RefType: JoinRoom,
			Field: "nickname",
			Reason: "Nickname should be 1 to 20 charecters long",
		})
	}

	joinedRoom, err := h.Cfg.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(joinedRoom)

	if joinedRoom.Started {
		return SendError(client, ErrorMessage{
			RefType: JoinRoom,
			Field: "roomLink",
			Reason: "Game already started",
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