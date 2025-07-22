package main

import (
	"log"
	"net/http"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/handlers"
	"server/internal/presence"
	"server/internal/store"
	"server/internal/transport"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := store.Init(redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisStore := store.New(redisClient)
	redisPresence := presence.NewRedisPresence(redisClient)
	localClients := domain.NewLocalClients()

	cfg := &config.Config{
		Environment:  config.Dev,
		Store:        redisStore,
		Presence:     redisPresence,
		LocalClients: localClients,
	}

	router := handlers.NewRouter(cfg)
	connectionManager := transport.NewConnectionManager(localClients, router)
	server := transport.NewServer(cfg, connectionManager)

	http.HandleFunc("/ws", server.HandleWebSocket)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
