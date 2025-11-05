package main

import (
	"context"
	"log"
	"net/http"
	"server/internal/broker"
	"server/internal/config"
	"server/internal/domain"
	"time"

	"server/internal/events"
	"server/internal/handlers"
	"server/internal/presence"
	"server/internal/store"
	"server/internal/transport"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// redisStore := store.NewRedisStore(redisClient)
	// redisPresence := presence.NewRedisPresence(redisClient)
	// redisBroker := broker.NewRedisBroker(redisClient)

	memoryStore := store.NewMemoryStore()
	memoryPresence := presence.NewMemoryPresence()
	memoryBroker := broker.NewMemoryBroker()

	localClients := domain.NewLocalClients()

	cfg := &config.Config{
		Environment:  config.Dev,
		// Store:        redisStore,
		Store:        memoryStore,
		// Presence:     redisPresence,
		Presence:     memoryPresence,
		// Broker: redisBroker,
		Broker: memoryBroker,
		LocalClients: localClients,
		DisconnectedClientTTL: time.Minute * 1,
	}

	eventHandler := events.NewRoomEventHandler(cfg)
	memoryBroker.SetEventCallback(eventHandler.HandleRoomEvent)

	router := handlers.NewRouter(cfg)
	connectionManager := transport.NewConnectionManager(cfg, router)
	server := transport.NewServer(cfg, connectionManager)

	http.HandleFunc("/ws", server.HandleWebSocket)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
