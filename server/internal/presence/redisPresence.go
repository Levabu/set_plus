package presence

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/room"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisPresence struct {
	client *redis.Client
}

func NewRedisPresence(client *redis.Client) *RedisPresence {
	return &RedisPresence{client: client}
}

func roomClientsKey(roomID uuid.UUID) string {
	return fmt.Sprintf("room:%s:client", roomID.String())
}

func roomChannel(roomID uuid.UUID) string {
	return fmt.Sprintf("room:%s:channel", roomID.String())
}

func (p *RedisPresence) JoinRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error {
	return p.client.SAdd(ctx, roomClientsKey(roomID), clientID.String()).Err()
}

func (p *RedisPresence) LeaveRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error {
	return p.client.SRem(ctx, roomClientsKey(roomID), clientID.String()).Err()
}

func (p *RedisPresence) GetRoomMembers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error) {
	members, err := p.client.SMembers(ctx, roomClientsKey(roomID)).Result()
	if err != nil {
		return nil, err
	}
	var ids []uuid.UUID
	for _, m := range members {
		id, err := uuid.Parse(m)
		if err == nil {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func (p *RedisPresence) SubscribeToRoom(ctx context.Context, roomID uuid.UUID, handler func(clientID uuid.UUID, event room.Event)) error {
	sub := p.client.Subscribe(ctx, roomChannel(roomID))
	ch := sub.Channel()

	log.Println("subscribed to channel:", roomChannel(roomID))
	go func() {
		for msg := range ch {
			var envelope struct {
				ClientID string          `json:"clientID"`
				Payload  room.Event `json:"payload"`
			}
			// log.Println("event received:", msg)
			if err := json.Unmarshal([]byte(msg.Payload), &envelope); err != nil {
				log.Println("invalid broadcast message:", err)
				continue
			}
			// log.Println("event unmarshalled:", envelope)
			clientID, err := uuid.Parse(envelope.ClientID)
			if err != nil {
				log.Println("invalid client ID in broadcast:", err)
				continue
			}
			log.Println("received room event:", envelope.Payload.Type)
			handler(clientID, envelope.Payload)
		}
	}()

	return nil
}