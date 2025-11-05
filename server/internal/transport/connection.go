package transport

import (
	"context"
	"encoding/json"
	"log"
	"server/internal/config"
	"server/internal/domain"
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

	cm.cfg.LocalClients.Add(client)
	defer cm.cfg.LocalClients.Remove(client.ID)

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

func (cm *ConnectionManager) CleanupRoom(roomID uuid.UUID) {
	cm.cfg.LocalClients.CleanupLocalRoomClients(roomID)
	cm.cfg.Presence.CleanupPresenceRoom(context.Background(), roomID)
	cm.cfg.Store.CleanupStoreRoom(context.Background(), roomID)
}
