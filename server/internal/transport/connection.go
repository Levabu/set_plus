package transport

import (
	"encoding/json"
	"log"
	"server/internal/domain"
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

