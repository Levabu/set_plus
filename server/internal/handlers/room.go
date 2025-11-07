package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/events"

	"github.com/google/uuid"
)

type RoomHandler struct {
	config       *config.Config
	eventHandler *events.RoomEventHandler
}

func NewRoomHandler(cfg *config.Config) *RoomHandler {
	return &RoomHandler{
		config:       cfg,
		eventHandler: events.NewRoomEventHandler(cfg),
	}
}

func (h *RoomHandler) HandleCreateRoom(client *domain.LocalClient, rawMsg json.RawMessage) error {
	var msg domain.CreateRoomMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	if len(msg.Nickname) < 1 || len(msg.Nickname) > 20 {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.CreateRoom,
			Field:   "nickname",
			Reason:  "Nickname should be 1 to 20 characters long",
		})
	}
	client.Nickname = msg.Nickname

	newRoom := domain.Room{
		ID:      uuid.New(),
		OwnerID: client.ID,
		Started: false,
	}
	client.RoomID = newRoom.ID

	if err := h.config.Store.SetRoom(context.Background(), &newRoom); err != nil {
		return err
	}

	if err := h.config.Presence.JoinRoom(context.Background(), newRoom.ID, client.ID); err != nil {
		return err
	}

	go h.config.Broker.SubscribeToRoom(context.Background(), newRoom.ID, func(clientID uuid.UUID, msg any) {
		msgData, ok := msg.([]byte)
		if !ok {
			return
		}
		h.eventHandler.HandleRoomEventMessage(newRoom.ID, clientID, msgData)
	})

	domain.SendJSON(client, domain.CreatedRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.CreatedRoom},
		RoomID:         newRoom.ID,
		PlayerID:       newRoom.OwnerID,
		Nickname:       msg.Nickname,
	})
	return nil
}

func (h *RoomHandler) HandleJoinRoom(client *domain.LocalClient, rawMsg json.RawMessage) error {
	var msg domain.JoinRoomMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	if len(msg.Nickname) < 1 || len(msg.Nickname) > 20 {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.JoinRoom,
			Field:   "nickname",
			Reason:  "Nickname should be 1 to 20 characters long",
		})
	}
	client.Nickname = msg.Nickname

	joinedRoom, err := h.config.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.JoinRoom,
			Field:   "roomLink",
			Reason:  "Room doesn't exist",
		})
	}

	if joinedRoom.Started {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.JoinRoom,
			Field:   "roomLink",
			Reason:  "Game already started",
		})
	}
	client.RoomID = joinedRoom.ID

	if err := h.config.Presence.JoinRoom(context.Background(), joinedRoom.ID, client.ID); err != nil {
		return err
	}

	// Send response to the joining client first
	domain.SendJSON(client, domain.JoinedRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.JoinedRoom},
		RoomID:         joinedRoom.ID,
		PlayerID:       client.ID,
		Nickname:       msg.Nickname,
	})

	// Send existing room members to the newly joined player
	existingMembers, err := h.config.Presence.GetActiveRoomMembersIDs(context.Background(), joinedRoom.ID)
	if err == nil {
		for _, memberID := range existingMembers {
			// Don't send the player their own join message again
			if memberID == client.ID {
				continue
			}
			memberClient := h.config.LocalClients.Get(memberID)
			if memberClient != nil {
				// Send info about existing member to the new player
				domain.SendJSON(client, domain.JoinedRoomMessage{
					BaseOutMessage: domain.BaseOutMessage{Type: domain.JoinedRoom},
					RoomID:         joinedRoom.ID,
					PlayerID:       memberID,
					Nickname:       memberClient.Nickname,
				})
			}
		}
	}

	log.Printf("Player %s (%s) joined room %s", client.ID, client.Nickname, joinedRoom.ID)

	// Publish room event to notify other members
	err = h.config.Broker.PublishRoomUpdate(context.Background(), joinedRoom.ID, domain.Event{
		Type:     domain.PlayerJoinedEvent,
		CliendID: client.ID,
	})
	if err != nil {
		return err
	}

	return nil
}
