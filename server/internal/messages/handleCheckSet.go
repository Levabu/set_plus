package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/game"
	"server/internal/room"
	"server/internal/server"

	"github.com/google/uuid"
)

func (h *Handler) handleCheckSet(client *server.Client, rawMsg json.RawMessage) error {
	var msg CheckSetMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	// validate
	r, err := h.Cfg.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return err
	}

	if msg.GameID != r.GameID {
		return SendError(client, ErrorMessage{
			RefType: CheckSet,
			Reason:  "Incorrect game id",
		})
	}

	gameState, err := h.Cfg.Store.GetGameState(context.Background(), r.GameID)
	if err != nil {
		return err
	}

	if gameState.Finished {
		return SendError(client, ErrorMessage{
			RefType: CheckSet,
			Reason:  "game already finished",
		})
	}

	if err := validateSetInput(gameState, msg.CardIDs); err != nil {
		return SendError(client, ErrorMessage{
			RefType: CheckSet,
			Reason:  err.Error(),
		})
	}

	// handle
	cards := make([]game.Card, len(msg.CardIDs))
	for i, id := range msg.CardIDs {
		card, ok := (*gameState.Cards)[id]
		if !ok {
			return fmt.Errorf("card missing after validation")
		}
		cards[i] = card
	}

	isSet := gameState.IsSet(cards)
	if !isSet {
		SendJSON(client, CheckSetResultMessage{
			BaseOutMessage: BaseOutMessage{Type: CheckSetResult},
			IsSet:          false,
		})
		return nil
	}

	SendJSON(client, CheckSetResultMessage{
		BaseOutMessage: BaseOutMessage{Type: CheckSetResult},
		IsSet:          true,
	})

	gameState.DiscardCards(cards)
	gameState.DealCards(gameState.GameConfig.VariationsNumber)
	gameState.DealCardsUntilSetAvailable(gameState.GameConfig.VariationsNumber, 30)

	player := (*gameState.Players)[client.ID]
	player.Score += 1
	(*gameState.Players)[client.ID] = player

	gameOver := gameState.IsGameOver()
	gameState.Finished = gameOver

	err = h.Cfg.Store.SetGameState(context.Background(), gameState)
	if err != nil {
		return err
	}

	var eventType room.EventType
	if gameOver {
		eventType = room.GameOver
	} else {
		eventType = room.ChangedGameState
	}

	h.Cfg.Store.PublishRoomUpdate(context.Background(), r.ID, room.Event{
		Type: eventType,
		CliendID: client.ID,
	})

	return nil
}

func validateSetInput(game *game.Game, ids []uuid.UUID) error {
	cardsSet := make(map[uuid.UUID]struct{})
	for _, id := range ids {
		cardsSet[id] = struct{}{}

		card, ok := (*game.Cards)[id]
		if !ok || !card.IsVisible || card.IsDiscarded {
			return fmt.Errorf("card not in play")
		}
	}
	if len(cardsSet) != game.GameConfig.VariationsNumber {
		return fmt.Errorf("duplicate card")
	}
	return nil
}
