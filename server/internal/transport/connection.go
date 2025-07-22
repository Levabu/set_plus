package transport

import (
	"encoding/json"
	"log"
	"server/internal/domain"

	"github.com/google/uuid"
)

type ConnectionManager struct {
	clients domain.LocalClientManager
	router  domain.MessageRouter
}

func NewConnectionManager(clients domain.LocalClientManager, router domain.MessageRouter) *ConnectionManager {
	return &ConnectionManager{
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
			log.Println("Read error:", err)
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

func (cm *ConnectionManager) BroadcastToRoom(roomID uuid.UUID, message interface{}) error {
	allClients := cm.clients.GetAll()
	for _, client := range allClients {
		if client.RoomID == roomID {
			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("Error broadcasting to client %s: %v", client.ID, err)
			}
		}
	}
	return nil
}

func (cm *ConnectionManager) SendToClient(clientID uuid.UUID, message interface{}) error {
	client := cm.clients.Get(clientID)
	if client == nil {
		return nil
	}
	return client.Conn.WriteJSON(message)
}
