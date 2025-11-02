package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"server/internal/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisBroker struct {
	client *redis.Client
}

func NewRedisBroker(client *redis.Client) *RedisBroker {
	return &RedisBroker{
		client: client,
	}
}

func roomChannel(roomID uuid.UUID) string {
	return fmt.Sprintf("room:%s:channel", roomID.String())
}

func (p *RedisBroker) SubscribeToRoom(ctx context.Context, roomID uuid.UUID, handler func(clientID uuid.UUID, msg any)) error {
	sub := p.client.Subscribe(ctx, roomChannel(roomID))
	ch := sub.Channel()

	go func() {
		for msg := range ch {
			var envelope struct {
				ClientID string          `json:"clientID"`
				Payload  json.RawMessage `json:"payload"`
			}
			if err := json.Unmarshal([]byte(msg.Payload), &envelope); err != nil {
				log.Println("invalid broadcast message:", err)
				continue
			}
			clientID, err := uuid.Parse(envelope.ClientID)
			if err != nil {
				log.Println("invalid client ID in broadcast:", err)
				continue
			}
			handler(clientID, []byte(envelope.Payload))
		}
	}()

	return nil
}

func (s *RedisBroker) PublishRoomUpdate(ctx context.Context, roomID uuid.UUID, event domain.Event) error {
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
