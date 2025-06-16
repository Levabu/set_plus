package messages

import (
	"context"
	"log"

	"server/internal/room"

	"github.com/google/uuid"
)

func (h *Handler) BroadcastToRoom(ctx context.Context, roomID uuid.UUID, payload interface{}) error {
	// log.Println("getting room members")
	players, err := h.Cfg.Presence.GetRoomMembers(ctx, roomID)
	if err != nil {
		return err
	}
	// log.Println("room members:", cliendsIDs)

	for _, player := range players {
		// log.Println(h.Cfg.LocalClients)
		if h.Cfg.LocalClients == nil {
			continue
		}
		client := h.Cfg.LocalClients.Get(player.ID)
		if client == nil {
			continue
		}
		SendJSON(client, payload)
	}

	return nil
}

func (h *Handler) HandleRoomEvent(id uuid.UUID, event room.Event) {
	switch event.Type {
	case room.JoinedPlayer:
		h.BroadcastToRoom(context.Background(), id, JoinedRoomMessage{
			BaseOutMessage: BaseOutMessage{Type: JoinedRoom},
			RoomID:         id,
			PlayerID:       event.CliendID,
		})
	case room.StartedGame:
		log.Println("handling event: started game")
		room, err := h.Cfg.Store.GetRoom(context.Background(), id)
		if err != nil {
			return
		}
		log.Println("starteg game in room:", room.ID)
		game, err := h.Cfg.Store.GetGameState(context.Background(), room.GameID)
		if err != nil {
			return
		}
		h.BroadcastToRoom(context.Background(), id, StartedGameMessage{
			BaseOutMessage: BaseOutMessage{Type: StartedGame},
			GameID:         room.GameID,
			GameVersion:    game.GameVersion,
			Deck:           game.GetInPlayCards(),
			Players:        *game.Players,
		})
	case room.ChangedGameState:
		room, err := h.Cfg.Store.GetRoom(context.Background(), id)
		if err != nil {
			return
		}
		game, err := h.Cfg.Store.GetGameState(context.Background(), room.GameID)
		if err != nil {
			return
		}
		h.BroadcastToRoom(context.Background(), id, ChangedGameStateMessage{
			BaseOutMessage: BaseOutMessage{Type: ChangedGameState},
			GameID:         room.GameID,
			Deck:           game.GetInPlayCards(),
			Players:        *game.Players,
		})
	case room.GameOver:
		room, err := h.Cfg.Store.GetRoom(context.Background(), id)
		if err != nil {
			return
		}
		game, err := h.Cfg.Store.GetGameState(context.Background(), room.GameID)
		if err != nil {
			return
		}
		h.BroadcastToRoom(context.Background(), id, GameOverMessage{
			BaseOutMessage: BaseOutMessage{Type: GameOver},
			GameID:         game.GameID,
			Deck:           game.GetInPlayCards(),
			Players:        *game.Players,
		})
	}
}
