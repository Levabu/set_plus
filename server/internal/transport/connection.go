package transport

import (
	"context"
	"encoding/json"
	"log"
	"server/internal/config"
	"server/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	cfg     *config.Config
	clients domain.LocalClientManager
	router  domain.MessageRouter
}

func NewConnectionManager(cfg *config.Config, clients domain.LocalClientManager, router domain.MessageRouter) *ConnectionManager {
	return &ConnectionManager{
		cfg:     cfg,
		clients: clients,
		router:  router,
	}
}

func (cm *ConnectionManager) HandleConnection(client *domain.LocalClient) {
	defer client.Conn.Close()

	cm.clients.Add(client)
	defer cm.clients.Remove(client.ID)

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
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
	cm.clients.SetClientConnected(client.ID, false)
	cm.cfg.Presence.LeaveRoom(context.Background(), client.ID)

	err := cm.cfg.Broker.PublishRoomUpdate(context.Background(), client.RoomID, domain.Event{
		Type:     domain.PlayerLeftEvent,
		CliendID: client.ID,
	})

	// set timer for local cleanup
	return err
}

func (cm *ConnectionManager) CleanupRoom(roomID uuid.UUID) {
	cm.clients.CleanupLocalRoomClients(roomID)
	cm.cfg.Presence.CleanupPresenceRoom(context.Background(), roomID)
	cm.cfg.Store.CleanupStoreRoom(context.Background(), roomID)
}
