package transport

import (
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/domain"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Server struct {
	upgrader          websocket.Upgrader
	config            *config.Config
	connectionManager domain.ConnectionManager
}

func NewServer(cfg *config.Config, connectionManager domain.ConnectionManager) *Server {
	return &Server{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
		config:            cfg,
		connectionManager: connectionManager,
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	queryParams := r.URL.Query()
	clientID, err := uuid.Parse(queryParams.Get("clientID"))

	var client *domain.LocalClient
	client = s.config.LocalClients.Get(clientID)
	if err == nil && clientID != uuid.Nil && client != nil && client.RoomID != uuid.Nil {
		if client.Conn != nil {
			client.Conn.Close()
		}
		client.Conn = conn
		client.WriteChan = make(chan interface{}, 256)
	} else {
		client = &domain.LocalClient{
			ID:   uuid.New(),
			Conn: conn,
			Connected: true,
			WriteChan: make(chan interface{}, 256),
		}
	}

	go s.connectionManager.HandleConnection(client)
}
