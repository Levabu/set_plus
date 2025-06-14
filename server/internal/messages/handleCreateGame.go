package messages

import (
	"encoding/json"
	"fmt"
	"log"

	"server/internal/game"
	"server/internal/server"

	"github.com/google/uuid"
)

func handleCreateGame(client *server.Client, rawMsg json.RawMessage) error {
	type Response struct {
		BaseMessage
		GameID uuid.UUID   `json:"gameID"`
		Deck   []game.Card `json:"deck"`
	}

	var msg CreateGameMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	if !msg.GameVersion.IsValid() {
		return fmt.Errorf("unsupported game version: %s", msg.GameVersion)
	}

	gameConfig, exists := game.GameVersions[msg.GameVersion]
	if !exists {
		return fmt.Errorf("unsupported game version: %s", msg.GameVersion)
	}

	game, err := game.NewGame(gameConfig)
	if err != nil {
		return fmt.Errorf("unable to create new game: %s", err.Error())
	}
	game.GenerateCards()
	game.ShuffleDeck()
	game.DealCards(gameConfig.InitialDeal)
	for {
		if game.IsSetAvailable() {
			break
		}
		game.DealCards(gameConfig.VariationsNumber)
	}

	sendJSON(client, Response{
		BaseMessage: BaseMessage{Type: "GAME_CREATED"},
		GameID: uuid.New(),
		Deck: game.Deck,
	})

	log.Println("send create game message")

	return nil
}
