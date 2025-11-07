package events

import (
	"context"
	// "log"
	"server/internal/domain"

	"github.com/google/uuid"
)

func (h *RoomEventHandler) handleJoinedPlayer(roomID uuid.UUID, event domain.Event) error {
	// Get the nickname of the player who joined
	joinedClient := h.config.LocalClients.Get(event.CliendID)
	nickname := ""
	if joinedClient != nil {
		nickname = joinedClient.Nickname
		// log.Printf("Found joined client %s with nickname: %s", event.CliendID, nickname)
	} else {
		// log.Printf("Could not find joined client %s in LocalClients", event.CliendID)
	}

	message := domain.JoinedRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.JoinedRoom},
		RoomID:         roomID,
		PlayerID:       event.CliendID,
		Nickname:       nickname,
	}

	// Get all room members and broadcast to everyone EXCEPT the player who just joined
	// (they already got their own JoinedRoom message from the handler)
	members, err := h.config.Presence.GetActiveRoomMembersIDs(context.Background(), roomID)
	if err != nil {
		return err
	}

	for _, memberID := range members {
		// Don't send to the player who just joined - they already got their response
		if memberID == event.CliendID {
			continue
		}
		memberClient := h.config.LocalClients.Get(memberID)
		if memberClient != nil {
			// log.Printf("Sending joined player notification to %s", memberID)
			domain.SendJSON(memberClient, message)
		} else {
			// log.Printf("Could not find member client %s in LocalClients", memberID)
		}
	}
	return nil
}

func (h *RoomEventHandler) handleDisconnectedPlayer(roomID uuid.UUID, event domain.Event) error {
	h.config.LocalClients.SetClientConnected(event.CliendID, false)

	// set timer for local cleanup

	msg := domain.LeftRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.LeftRoom},
		PlayerID:       event.CliendID,
	}

	return h.BroadcastToRoom(context.Background(), roomID, msg, h.config.LocalClients)
}

func (h *RoomEventHandler) handleReconnectedPlayer(roomID uuid.UUID, event domain.Event) error {
	h.config.LocalClients.SetClientConnected(event.CliendID, true)

	msg := domain.ReconnectedToRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.LeftRoom},
		PlayerID:       event.CliendID,
	}

	return h.BroadcastToRoom(context.Background(), roomID, msg, h.config.LocalClients)
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

	return h.BroadcastToRoom(context.Background(), roomID, startedMessage, h.config.LocalClients)
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

	return h.BroadcastToRoom(context.Background(), roomID, changedMessage, h.config.LocalClients)
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

	return h.BroadcastToRoom(context.Background(), roomID, gameOverMessage, h.config.LocalClients)
}
