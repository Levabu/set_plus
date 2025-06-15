package messages

import (
	"context"
	"log"

	"server/internal/room"

	"github.com/google/uuid"
)

func (h *Handler) BroadcastToRoom(ctx context.Context, roomID uuid.UUID, payload interface{}) error {
	// log.Println("getting room members")
	cliendsIDs, err := h.Cfg.Presence.GetRoomMembers(ctx, roomID)
	if err != nil {
		return err
	}
	// log.Println("room members:", cliendsIDs)

	for _, clientID := range cliendsIDs {
		// log.Println(h.Cfg.LocalClients)
		if h.Cfg.LocalClients == nil {
			continue
		}
		client := h.Cfg.LocalClients.Get(clientID)
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
		// log.Println("game state on start:", game)
		if err != nil {
			return
		}
		log.Println(room, event)
		h.BroadcastToRoom(context.Background(), id, StartedGameMessage{
			BaseOutMessage: BaseOutMessage{Type: StartedGame},
			GameID:         room.GameID,
			Deck:           game.Deck,
		})
	case room.ChangedGameState:
		room, err := h.Cfg.Store.GetRoom(context.Background(), id)
		if err != nil {
			return
		}
		log.Println("changed game in room:", room.ID)
		game, err := h.Cfg.Store.GetGameState(context.Background(), room.GameID)
		if err != nil {
			return
		}
		log.Println(room, event)
		h.BroadcastToRoom(context.Background(), id, ChangedGameStateMessage{
			BaseOutMessage: BaseOutMessage{Type: ChangedGameState},
			GameID:         room.GameID,
			Deck:           game.Deck,
		})
	}
}
