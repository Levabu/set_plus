package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/game"
	"server/internal/room"
	"server/internal/server"
)

func (h *Handler) handleStartGame(client *server.Client, rawMsg json.RawMessage) error {
	var msg StartGameMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	game, err := createNewGame(msg)
	if err != nil {
		return err
	}

	r, err := h.Cfg.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return err
	}
	r.GameID = game.GameID

	if r.OwnerID != client.ID {
		return SendError(client, ErrorMessage{
			RefType: StartGame,
			Reason:  "only owner of the room can start the game",
		})
	}

	r.Started = true

	if err = h.Cfg.Store.SetRoom(context.Background(), r); err != nil {
		return err
	}

	players, err := h.Cfg.Presence.GetRoomMembers(context.Background(), r.ID)
	if err != nil {
		return err
	}
	for _, player := range players {
		player.Score = 0
		(*game.Players)[player.ID] = player
	}

	if err = h.Cfg.Store.SetGameState(context.Background(), game); err != nil {
		return err
	}

	err = h.Cfg.Store.PublishRoomUpdate(context.Background(), client.RoomID, room.Event{
		Type:     room.StartedGame,
		CliendID: client.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func createNewGame(msg StartGameMessage) (*game.Game, error) {
	if !msg.GameVersion.IsValid() {
		return nil, fmt.Errorf("unsupported game version: %s", msg.GameVersion)
	}

	game, err := game.NewGame(msg.GameVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to create new game: %s", err.Error())
	}
	game.GenerateCards()
	game.ShuffleDeck()
	game.DealCards(game.GameConfig.InitialDeal)
	game.DealCardsUntilSetAvailable(game.GameConfig.VariationsNumber, 30)

	return game, nil
}
