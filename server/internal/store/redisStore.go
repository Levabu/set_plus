package store

import (
	"context"
	"encoding/json"
	"fmt"
	"server/internal/domain"
	"server/internal/game"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
	}
}

func (s *RedisStore) SetRoom(ctx context.Context, room *domain.Room) error {
	key := fmt.Sprintf("room:%s", room.ID)
	data, err := json.Marshal(room)
	if err != nil {
		return fmt.Errorf("error setting a room: %s", err)
	}
	return s.client.Set(ctx, key, data, 0).Err()
}

func (s *RedisStore) GetRoom(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	key := fmt.Sprintf("room:%s", id)
	data, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var room domain.Room
	if err := json.Unmarshal([]byte(data), &room); err != nil {
		return nil, err
	}
	return &room, nil
}

func (s *RedisStore) SetGameState(ctx context.Context, game *game.Game) error {
	key := fmt.Sprintf("game:%s", game.GameID)
	data, err := json.Marshal(game)
	if err != nil {
		return fmt.Errorf("error saving game: %s", err)
	}
	return s.client.Set(ctx, key, data, 0).Err()
}

func (s *RedisStore) GetGameState(ctx context.Context, id uuid.UUID) (*game.Game, error) {
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

func (s *RedisStore) PublishRoomUpdate(ctx context.Context, roomID uuid.UUID, event domain.Event) error {
	envelope := struct {
		ClientID string       `json:"clientID"`
		Payload  domain.Event `json:"payload"`
	}{
		ClientID: event.CliendID.String(),
		Payload:  event,
	}
	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	channel := fmt.Sprintf("room:%s:channel", roomID)
	return s.client.Publish(ctx, channel, data).Err()
}
