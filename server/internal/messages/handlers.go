package messages

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"server/internal/game"
	"server/internal/room"
	"server/internal/server"

	"github.com/google/uuid"
)

func (h *Handler) handleCreateRoom(client *server.Client, rawMsg json.RawMessage) error {
	newRoom := room.Room{
		ID:      uuid.New(),
		OwnerID: client.ID,
		Started: false,
	}
	client.RoomID = newRoom.ID

	if err := h.Cfg.Store.SetRoom(context.Background(), &newRoom); err != nil {
		return err
	}

	if err := h.Cfg.Presence.JoinRoom(context.Background(), newRoom.ID, client.ID); err != nil {
		return err
	}

	go h.Cfg.Presence.SubscribeToRoom(context.Background(), newRoom.ID, func(clientID uuid.UUID, event room.Event) {
		h.HandleRoomEvent(clientID, event)
	})

	SendJSON(client, CreatedRoomMessage{
		BaseOutMessage: BaseOutMessage{Type: CreatedRoom},
		RoomID:         newRoom.ID,
		PlayerID:       newRoom.OwnerID,
	})
	return nil
}

func (h *Handler) handleJoinRoom(client *server.Client, rawMsg json.RawMessage) error {
	var msg JoinRoomMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	joinedRoom, err := h.Cfg.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return err
	}

	res := JoinedRoomMessage{
		BaseOutMessage: BaseOutMessage{Type: JoinedRoom},
		RoomID:         joinedRoom.ID,
	}

	if joinedRoom.Started {
		res.Error = "game already started"
		SendJSON(client, res)
	}

	if err := h.Cfg.Presence.JoinRoom(context.Background(), joinedRoom.ID, client.ID); err != nil {
		return err
	}

	err = h.Cfg.Store.PublishRoomUpdate(context.Background(), joinedRoom.ID, room.Event{
		Type:     room.JoinedPlayer,
		CliendID: client.ID,
	})
	if err != nil {
		return err
	}

	SendJSON(client, res)
	return nil
}

func (h *Handler) handleStartGame(client *server.Client, rawMsg json.RawMessage) error {
	var msg StartGameMessage
	if err := json.Unmarshal(rawMsg, &msg); err != nil {
		return fmt.Errorf("invalid message: %s", err.Error())
	}

	game, err := _createNewGame(msg)
	if err != nil {
		return err
	}

	r, err := h.Cfg.Store.GetRoom(context.Background(), msg.RoomID)
	if err != nil {
		return err
	}
	r.GameID = game.GameID

	// todo: only allow owner to start
	// if r.OwnerID != client.ID

	if err = h.Cfg.Store.SetRoom(context.Background(), r); err != nil {
		return err
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

func _createNewGame(msg StartGameMessage) (*game.Game, error) {
	if !msg.GameVersion.IsValid() {
		return nil, fmt.Errorf("unsupported game version: %s", msg.GameVersion)
	}

	gameConfig, exists := game.GameVersions[msg.GameVersion]
	if !exists {
		return nil, fmt.Errorf("unsupported game version: %s", msg.GameVersion)
	}

	game, err := game.NewGame(gameConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create new game: %s", err.Error())
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

	return game, nil
}

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
	}
}
