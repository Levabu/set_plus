package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/game"

	// "server/internal/presence"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	cfg    *config.Config
	router domain.MessageRouter
}

func NewConnectionManager(cfg *config.Config, router domain.MessageRouter) *ConnectionManager {
	return &ConnectionManager{
		cfg:    cfg,
		router: router,
	}
}

func (cm *ConnectionManager) HandleConnection(client *domain.LocalClient) {
	defer client.Conn.Close()
	defer close(client.WriteChan)

	if client.Connected {
		cm.cfg.LocalClients.Add(client)
	} else {
		cm.HandleReconnection(client)
	}

	go cm.StartWriter(client)
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Println("Websocket error:", err)
			}
			cm.HandleDisconnection(client)
			break
		}

		var baseMessage domain.InMessage
		if err := json.Unmarshal(msg, &baseMessage); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		if err := cm.router.HandleMessage(client, baseMessage.Type, msg); err != nil {
			log.Printf("error handling message %s: %v", baseMessage.Type, err)
		}
	}
}

func (cm *ConnectionManager) HandleDisconnection(client *domain.LocalClient) error {
	clientID := client.ID
	roomID := client.RoomID

	cm.cfg.LocalClients.SetClientConnected(clientID, false)
	cm.cfg.Presence.LeaveRoom(context.Background(), clientID)

	var err error
	if roomID != uuid.Nil {
		err = cm.cfg.Broker.PublishRoomUpdate(context.Background(), roomID, domain.Event{
			Type:     domain.PlayerLeftEvent,
			CliendID: clientID,
		})
	}

	// set timer for local cleanup
	if client.ReconnectTimer != nil {
		client.ReconnectTimer.Stop()
	}
	client.ReconnectTimer = time.AfterFunc(cm.cfg.DisconnectedClientTTL, func() {
		time.Sleep(cm.cfg.DisconnectedClientTTL)

		client := cm.cfg.LocalClients.Get(clientID)
		if client == nil || client.Connected || time.Since(client.DisconnectedAt) < cm.cfg.DisconnectedClientTTL {
			return
		}

		cm.cfg.LocalClients.Remove(clientID)
		cm.cfg.Presence.RemoveClient(context.Background(), clientID, roomID)
		if roomID != uuid.Nil && cm.cfg.LocalClients.IsRoomEmpty(roomID) {
			cm.CleanupRoom(roomID)
		}
	})

	return err
}

func (cm *ConnectionManager) HandleReconnection(client *domain.LocalClient) error {
	fmt.Printf("Reconnect: %s\n", client.ID)
	if client.ReconnectTimer != nil {
		client.ReconnectTimer.Stop()
	}
	// update locally
	cm.cfg.LocalClients.SetClientConnected(client.ID, true)
	// notify
	cm.cfg.Presence.JoinRoom(context.Background(), client.RoomID, client.ID)
	cm.cfg.Broker.PublishRoomUpdate(context.Background(), client.RoomID, domain.Event{
		Type:     domain.PlayerReconnectedEvent,
		CliendID: client.ID,
	})
	// send current state to the player
	msg := domain.SendStateToReconnectedMessage{BaseOutMessage: domain.BaseOutMessage{Type: domain.SendStateToReconnected}}

	room, err := cm.cfg.Store.GetRoom(context.Background(), client.RoomID)
	if err != nil || room == nil {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.ReconnectToRoom,
			Reason:  "Room doesn't exist",
		})
	}
	msg.IsOwner = room.OwnerID == client.ID
	msg.RoomID = room.ID
	msg.Started = room.Started

	if room.Started && room.GameID != uuid.Nil {
		msg.GameID = room.GameID
		game, err := cm.cfg.Store.GetGameState(context.Background(), room.GameID)
		if err != nil || game == nil {
			return domain.SendError(client, domain.ErrorMessage{
				RefType: domain.ReconnectToRoom,
				Reason:  "Game doesn't exist",
			})
		}
		msg.GameVersion = game.GameVersion
		msg.Deck = game.Deck
		msg.Players = *game.Players
	} else {
		activeClients, err := cm.cfg.Presence.GetActiveRoomMembers(context.Background(), room.ID)
		if err != nil {
			return err
		}
		players := make(map[uuid.UUID]game.Player, 0)
		for _, c := range activeClients {
			if c.ID == client.ID {
				continue
			}
			players[c.ID] = game.Player{ID: c.ID}
		}
		msg.Players = players
	}
	fmt.Printf("SENDING STATE: %v\n", msg.IsOwner)
	bytes, err := json.Marshal(msg)
	fmt.Println(err)
	fmt.Println(string(bytes))
	return domain.SendJSON(client, msg)
}

func (cm *ConnectionManager) CleanupRoom(roomID uuid.UUID) {
	cm.cfg.LocalClients.CleanupLocalRoomClients(roomID)
	cm.cfg.Presence.CleanupPresenceRoom(context.Background(), roomID)
	cm.cfg.Store.CleanupStoreRoom(context.Background(), roomID)
}

func (cm *ConnectionManager) StartWriter(client *domain.LocalClient) {
	for msg := range client.WriteChan {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}
}