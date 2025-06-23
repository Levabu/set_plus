package events

import (
	"context"
	"encoding/json"
	"log"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/room"

	"github.com/google/uuid"
)

type Publisher struct {
	connectionManager domain.ConnectionManager
}

func NewPublisher(connectionManager domain.ConnectionManager) *Publisher {
	return &Publisher{
		connectionManager: connectionManager,
	}
}

func (p *Publisher) PublishRoomEvent(ctx context.Context, roomID uuid.UUID, event interface{}) error {
	return p.connectionManager.BroadcastToRoom(roomID, event)
}

func (p *Publisher) PublishGameEvent(ctx context.Context, gameID uuid.UUID, event interface{}) error {
	return nil
}

type RoomEventHandler struct {
	config *config.Config
}

func NewRoomEventHandler(config *config.Config) *RoomEventHandler {
	return &RoomEventHandler{
		config: config,
	}
}

// This is called from Redis PubSub subscription
func (h *RoomEventHandler) HandleRoomEventMessage(roomID uuid.UUID, clientID uuid.UUID, msgData []byte) {
	var event room.Event
	if err := json.Unmarshal(msgData, &event); err != nil {
		log.Printf("Failed to unmarshal room event: %v", err)
		return
	}

	if err := h.HandleRoomEvent(roomID, event); err != nil {
		log.Printf("Failed to handle room event: %v", err)
	}
}

func (h *RoomEventHandler) HandleRoomEvent(roomID uuid.UUID, event room.Event) error {
	switch event.Type {
	case room.JoinedPlayer:
		return h.handleJoinedPlayer(roomID, event)
	case room.StartedGame:
		return h.handleStartedGame(roomID, event)
	case room.ChangedGameState:
		return h.handleChangedGameState(roomID, event)
	case room.GameOver:
		return h.handleGameOver(roomID, event)
	}
	return nil
}

func (h *RoomEventHandler) handleJoinedPlayer(roomID uuid.UUID, event room.Event) error {
	// Get the nickname of the player who joined
	joinedClient := h.config.LocalClients.Get(event.CliendID)
	nickname := ""
	if joinedClient != nil {
		nickname = joinedClient.Nickname
		log.Printf("Found joined client %s with nickname: %s", event.CliendID, nickname)
	} else {
		log.Printf("Could not find joined client %s in LocalClients", event.CliendID)
	}

	message := domain.JoinedRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.JoinedRoom},
		RoomID:         roomID,
		PlayerID:       event.CliendID,
		Nickname:       nickname,
	}
	
	// Get all room members and broadcast to everyone EXCEPT the player who just joined
	// (they already got their own JoinedRoom message from the handler)
	members, err := h.config.Presence.GetRoomMembers(context.Background(), roomID)
	if err != nil {
		return err
	}

	log.Printf("Broadcasting joined player %s (%s) to %d members", event.CliendID, nickname, len(members))

	for _, memberID := range members {
		// Don't send to the player who just joined - they already got their response
		if memberID == event.CliendID {
			continue
		}
		memberClient := h.config.LocalClients.Get(memberID)
		if memberClient != nil {
			log.Printf("Sending joined player notification to %s", memberID)
			memberClient.Conn.WriteJSON(message)
		} else {
			log.Printf("Could not find member client %s in LocalClients", memberID)
		}
	}
	return nil
}

func (h *RoomEventHandler) handleStartedGame(roomID uuid.UUID, event room.Event) error {
	// Get room and game state
	gameRoom, err := h.config.Store.GetRoom(context.Background(), roomID)
	if err != nil {
		return err
	}
	
	gameState, err := h.config.Store.GetGameState(context.Background(), gameRoom.GameID)
	if err != nil {
		return err
	}

	startedMessage := domain.StartedGameMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.StartedGame},
		GameID:         gameState.GameID,
		GameVersion:    gameState.GameVersion,
		Deck:           gameState.GetVisibleCards(),
		Players:        *gameState.Players,
	}

	// Get all room members and broadcast
	members, err := h.config.Presence.GetRoomMembers(context.Background(), roomID)
	if err != nil {
		return err
	}

	for _, memberID := range members {
		memberClient := h.config.LocalClients.Get(memberID)
		if memberClient != nil {
			memberClient.Conn.WriteJSON(startedMessage)
		}
	}
	return nil
}

func (h *RoomEventHandler) handleChangedGameState(roomID uuid.UUID, event room.Event) error {
	// Get room and game state
	gameRoom, err := h.config.Store.GetRoom(context.Background(), roomID)
	if err != nil {
		return err
	}
	
	gameState, err := h.config.Store.GetGameState(context.Background(), gameRoom.GameID)
	if err != nil {
		return err
	}

	changedMessage := domain.ChangedGameStateMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.ChangedGameState},
		GameID:         gameState.GameID,
		Deck:           gameState.GetVisibleCards(),
		Players:        *gameState.Players,
	}

	// Get all room members and broadcast
	members, err := h.config.Presence.GetRoomMembers(context.Background(), roomID)
	if err != nil {
		return err
	}

	for _, memberID := range members {
		memberClient := h.config.LocalClients.Get(memberID)
		if memberClient != nil {
			memberClient.Conn.WriteJSON(changedMessage)
		}
	}
	return nil
}

func (h *RoomEventHandler) handleGameOver(roomID uuid.UUID, event room.Event) error {
	// Get room and game state
	gameRoom, err := h.config.Store.GetRoom(context.Background(), roomID)
	if err != nil {
		return err
	}
	
	gameState, err := h.config.Store.GetGameState(context.Background(), gameRoom.GameID)
	if err != nil {
		return err
	}

	gameOverMessage := domain.GameOverMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.GameOver},
		GameID:         gameState.GameID,
		Deck:           gameState.GetVisibleCards(),
		Players:        *gameState.Players,
	}

	// Get all room members and broadcast
	members, err := h.config.Presence.GetRoomMembers(context.Background(), roomID)
	if err != nil {
		return err
	}

	for _, memberID := range members {
		memberClient := h.config.LocalClients.Get(memberID)
		if memberClient != nil {
			memberClient.Conn.WriteJSON(gameOverMessage)
		}
	}
	return nil
}