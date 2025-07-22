package handlers

import (
	"encoding/json"
	"fmt"
	"server/internal/config"
	"server/internal/domain"
)

type Router struct {
	config   *config.Config
	handlers map[domain.InMessageType]domain.MessageHandler
}

func NewRouter(cfg *config.Config) *Router {
	r := &Router{
		config:   cfg,
		handlers: make(map[domain.InMessageType]domain.MessageHandler),
	}
	r.registerHandlers()
	return r
}

func (r *Router) registerHandlers() {
	roomHandler := NewRoomHandler(r.config)
	gameHandler := NewGameHandler(r.config)

	r.handlers = map[domain.InMessageType]domain.MessageHandler{
		domain.CreateRoom: roomHandler.HandleCreateRoom,
		domain.JoinRoom:   roomHandler.HandleJoinRoom,
		domain.StartGame:  gameHandler.HandleStartGame,
		domain.CheckSet:   gameHandler.HandleCheckSet,
	}
}

func (r *Router) RegisterHandler(msgType domain.InMessageType, handler domain.MessageHandler) {
	r.handlers[msgType] = handler
}

func (r *Router) HandleMessage(client *domain.LocalClient, msgType domain.InMessageType, rawMsg json.RawMessage) error {
	handler, ok := r.handlers[msgType]
	if !ok {
		return fmt.Errorf("unknown message type: %s", msgType)
	}
	return handler(client, rawMsg)
}
