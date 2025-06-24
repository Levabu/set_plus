package events

import (
	"context"
	"log"
	"server/internal/domain"

	"github.com/google/uuid"
)



func (h *RoomEventHandler) handleJoinedPlayer(roomID uuid.UUID, event domain.Event) error {
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

func (h *RoomEventHandler) handleStartedGame(roomID uuid.UUID, event domain.Event) error {
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

	// Use the new BroadcastToRoom method
	return h.config.Presence.BroadcastToRoom(context.Background(), roomID, startedMessage, h.config.LocalClients)
}

func (h *RoomEventHandler) handleChangedGameState(roomID uuid.UUID, event domain.Event) error {
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

	// Use the new BroadcastToRoom method
	return h.config.Presence.BroadcastToRoom(context.Background(), roomID, changedMessage, h.config.LocalClients)
}

func (h *RoomEventHandler) handleGameOver(roomID uuid.UUID, event domain.Event) error {
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

	// Use the new BroadcastToRoom method
	return h.config.Presence.BroadcastToRoom(context.Background(), roomID, gameOverMessage, h.config.LocalClients)
}


