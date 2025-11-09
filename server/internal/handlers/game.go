package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/config"
	"server/internal/domain"
	"server/internal/game"
	"time"

	"github.com/google/uuid"
)

type GameHandler struct {
	config *config.Config
}

func NewGameHandler(cfg *config.Config) *GameHandler {
	return &GameHandler{config: cfg}
}

func (h *GameHandler) HandleStartGame(client *domain.LocalClient, rawMsg json.RawMessage) error {
	var msg domain.StartGameMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	gameInstance, err := h.createNewGame(msg)
	if err != nil {
		return err
	}

	r, err := h.config.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return err
	}
	r.GameID = gameInstance.GameID

	if r.OwnerID != client.ID {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.StartGame,
			Reason:  "only owner of the room can start the game",
		})
	}

	r.Started = true

	if err = h.config.Store.SetRoom(context.Background(), r); err != nil {
		return err
	}

	players, err := h.config.Presence.GetActiveRoomMembersIDs(context.Background(), r.ID)
	if err != nil {
		return err
	}
	for _, playerID := range players {
		player := game.Player{
			ID:    playerID,
			Score: 0,
		}
		(*gameInstance.Players)[playerID] = player
	}

	if err = h.config.Store.SetGameState(context.Background(), gameInstance); err != nil {
		return err
	}

	err = h.config.Broker.PublishRoomUpdate(context.Background(), r.ID, domain.Event{
		Type:     domain.GameStartedEvent,
		CliendID: client.ID,
	})
	if err != nil {
		return err
	}
	return nil
}

func (h *GameHandler) HandleCheckSet(client *domain.LocalClient, rawMsg json.RawMessage) error {
	var msg domain.CheckSetMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	r, err := h.config.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return err
	}

	if msg.GameID != r.GameID {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.CheckSet,
			Reason:  "Incorrect game id",
		})
	}

	gameState, err := h.config.Store.GetGameState(context.Background(), r.GameID)
	if err != nil {
		return err
	}

	if gameState.Finished {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.CheckSet,
			Reason:  "game already finished",
		})
	}

	if err := h.validateSetInput(gameState, msg.CardIDs); err != nil {
		return domain.SendError(client, domain.ErrorMessage{
			RefType: domain.CheckSet,
			Reason:  err.Error(),
		})
	}

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
		domain.SendJSON(client, domain.CheckSetResultMessage{
			BaseOutMessage: domain.BaseOutMessage{Type: domain.CheckSetResult},
			IsSet:          false,
		})
		return nil
	}

	domain.SendJSON(client, domain.CheckSetResultMessage{
		BaseOutMessage: domain.BaseOutMessage{Type: domain.CheckSetResult},
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

	err = h.config.Store.SetGameState(context.Background(), gameState)
	if err != nil {
		return err
	}

	var eventType domain.EventType
	if gameOver {
		eventType = domain.GameOverEvent
	} else {
		eventType = domain.GameStateChangedEvent
	}

	h.config.Broker.PublishRoomUpdate(context.Background(), r.ID, domain.Event{
		Type:     eventType,
		CliendID: client.ID,
	})

	if gameOver {
		go func() {
			time.Sleep(time.Second * 3)
			h.config.Store.CleanupAfterGame(context.Background(), gameState.GameID)
		} ()
	}

	return nil
}

func (h *GameHandler) createNewGame(msg domain.StartGameMessage) (*game.Game, error) {
	if !msg.GameVersion.IsValid() {
		return nil, fmt.Errorf("unsupported game version: %s", msg.GameVersion)
	}

	gameInstance, err := game.NewGame(msg.GameVersion)
	if err != nil {
		return nil, fmt.Errorf("unable to create new game: %s", err.Error())
	}
	gameInstance.GenerateCards()
	gameInstance.ShuffleDeck()
	gameInstance.DealCards(gameInstance.GameConfig.InitialDeal)
	gameInstance.DealCardsUntilSetAvailable(gameInstance.GameConfig.VariationsNumber, 30)

	return gameInstance, nil
}

func (h *GameHandler) validateSetInput(gameState *game.Game, ids []uuid.UUID) error {
	cardsSet := make(map[uuid.UUID]struct{})
	for _, id := range ids {
		cardsSet[id] = struct{}{}

		card, ok := (*gameState.Cards)[id]
		if !ok || !card.IsVisible || card.IsDiscarded {
			return fmt.Errorf("card not in play")
		}
	}
	if len(cardsSet) != gameState.GameConfig.VariationsNumber {
		return fmt.Errorf("duplicate card")
	}
	return nil
}
