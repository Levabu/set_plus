package events

import (
	"context"
	"server/internal/domain"

	"github.com/google/uuid"
)

func (h *RoomEventHandler) handleJoinedPlayer(roomID uuid.UUID, event domain.Event) error {
	data := event.Data
	var nickname string
	if data != nil {
		nickname = data["nickname"]
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
		if memberID == event.CliendID {
			continue
		}
		memberClient := h.config.LocalClients.Get(memberID)
		if memberClient != nil {
			domain.SendJSON(memberClient, message)
		}
	}
	return nil
}

func (h *RoomEventHandler) handleDisconnectedPlayer(roomID uuid.UUID, event domain.Event) error {
	msg := domain.LeftRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.LeftRoom},
		PlayerID:       event.CliendID,
	}

	return h.BroadcastToRoom(context.Background(), roomID, msg, h.config.LocalClients)
}

func (h *RoomEventHandler) handleReconnectedPlayer(roomID uuid.UUID, event domain.Event) error {
	msg := domain.ReconnectedToRoomMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.ReconnectedToRoom},
		RoomID: roomID,
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
