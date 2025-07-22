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

	client := &domain.LocalClient{
		ID:   uuid.New(),
		Conn: conn,
	}

	go s.connectionManager.HandleConnection(client)
}
