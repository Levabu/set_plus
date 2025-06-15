package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/messages"
	"server/internal/presence"
	"server/internal/server"
	"server/internal/store"
	"sync"

	// "github.com/google/uuid"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow all origins (dev only)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, cfg *config.Config) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &server.Client{
		ID:   uuid.New(),
		Conn: conn,
	}

	go handleConnetion(cfg, client)
}

func handleConnetion(cfg *config.Config, client *server.Client) {
	defer client.Conn.Close()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)

		var baseMessage messages.InMessage
		if err := json.Unmarshal(msg, &baseMessage); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		log.Println("message type:", baseMessage.Type)
		handler := &messages.Handler{Cfg: cfg}
		messageHandlers := handler.RegisterHandlers()
		handle, ok := messageHandlers[baseMessage.Type]
		if !ok {
			log.Println("unknown message type:", baseMessage.Type)
			continue
		}

		err = handle(client, msg)
		if err != nil {
			log.Println("error handling message:", err.Error())
		}
	}
}

func main() {
	redisClient := store.Init(redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisStore := store.New(redisClient)
	presence := presence.NewRedisPresence(redisClient)
	localClients := server.LocalClients{
		Clients: make(map[uuid.UUID]*server.Client),
		Mu: sync.RWMutex{},
	}
	cfg := &config.Config{
		Environment:  config.Dev,
		Store:        redisStore,
		Presence:     presence,
		LocalClients: &localClients,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(w, r, cfg)
	})
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
