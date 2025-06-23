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
	presenceService := presence.NewRedisPresence(redisClient)
	clientManager := domain.NewLocalClients()

	cfg := &config.Config{
		Environment:  config.Dev,
		Store:        redisStore,
		Presence:     presenceService,
		LocalClients: clientManager,
	}

	router := handlers.NewRouter(cfg)
	connectionManager := transport.NewConnectionManager(clientManager, router)
	server := transport.NewServer(cfg, connectionManager)

	http.HandleFunc("/ws", server.HandleWebSocket)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
