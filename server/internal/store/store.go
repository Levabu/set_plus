package store

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/game"
	"server/internal/room"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func New(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) CreateRoom(ctx context.Context, room *room.Room) error {
	key := fmt.Sprintf("room:%s", room.ID)
	data, err := json.Marshal(room)
	if err != nil {
		return fmt.Errorf("error creating a room: %s", err)
	}
	return s.client.Set(ctx, key, data, 0).Err()
}

func (s *Store) GetRoom(ctx context.Context, id uuid.UUID) (*room.Room, error) {
	key := fmt.Sprintf("room:%s", id)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var room room.Room
	if err := json.Unmarshal([]byte(data), &room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *Store) SaveGameState(ctx context.Context, game *game.Game) error {
	key := fmt.Sprintf("game:%s", game.GameID)
	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("error saving game: %s", err)
	}
	return s.client.Set(ctx, key, data, 0).Err()
}

func (s *Store) GetGameState(ctx context.Context, id uuid.UUID) (*game.Game, error) {
	key := fmt.Sprintf("game:%s", id)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var game game.Game
	if err := json.Unmarshal([]byte(data), &game); err != nil {
		return nil, err
	}
	return &game, nil
}

func (s *Store) PublishRoomUpdate(ctx context.Context, id uuid.UUID, event room.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	channel := fmt.Sprintf("room:%s:events", id)
	return s.client.Publish(ctx, channel, data).Err()
}

// func (s *Store) SubscribeToRoomEvents(id uuid.UUID) {
// 	channel := fmt.Sprintf("room:%s:events", id)
// 	pubsub := s.client.Subscribe(context.Background(), channel)
// 	ch := pubsub.Channel()
// 	for msg := range ch {
// 		fmt.Sprintln(msg.String(), channel)
// 		// This is triggered when PublishRoomUpdate is called
// 		// handler.OnRoomUpdate(roomID, msg.Payload)
// 	}
// }
