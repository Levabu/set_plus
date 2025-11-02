package presence

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	// map room id to set of client ids
	return fmt.Sprintf("room:%s:clients", roomID.String())
}

func clientKey(clientID uuid.UUID) string {
	// map client id to client data
	return fmt.Sprintf("client:%s:status", clientID.String())
}

func (p *RedisPresence) GetClient(ctx context.Context, clientID uuid.UUID) (PresenceClient, error) {
	var status PresenceClient

	data, err := p.client.Get(ctx, clientKey(clientID)).Result()
	if err != nil {
		if err == redis.Nil {
			return status, fmt.Errorf("client status not found")
		}
		return status, err
	}

	err = json.Unmarshal([]byte(data), &status)
	return status, err
}

func (p *RedisPresence) SetClient(ctx context.Context, clientID uuid.UUID, status PresenceClient) error {
	status.LastSeen = time.Now().Unix()
	statusData, err := json.Marshal(status)
	if err != nil {
		return err
	}

	ttl := 5 * time.Minute
	if !status.Connected {
		ttl = 1 * time.Hour // Keep disconnected status longer for potential reconnection
	}

	return p.client.Set(ctx, clientKey(status.ID), statusData, ttl).Err()
}

func (p *RedisPresence) IsClientConnected(ctx context.Context, clientID uuid.UUID) (bool, error) {
	client, err := p.GetClient(ctx, clientID)
	if err != nil {
		return false, nil // If status not found, consider disconnected
	}

	now := time.Now().Unix()
	if now-client.LastSeen > 300 { // 5 minutes
		return false, nil
	}

	return client.Connected, nil
}

func (p *RedisPresence) JoinRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error {
	// store client status
	if err := p.SetClient(ctx, clientID, PresenceClient{
		ID:        clientID,
		RoomID:    roomID,
		Connected: true,
		LastSeen:  time.Now().Unix(),
	}); err != nil {
		return err
	}

	// add client to room set
	value := fmt.Sprint(clientID.String())
	return p.client.SAdd(ctx, roomClientsKey(roomID), value).Err()
}

func (p *RedisPresence) LeaveRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error {
	// update status
	if err := p.SetClient(ctx, clientID, PresenceClient{
		ID:        clientID,
		Connected: false,
		LastSeen:  time.Now().Unix(),
	}); err != nil {
		log.Printf("Failed to update client status on leave: %v", err)
	}

	// remove from room
	return p.client.SRem(ctx, roomClientsKey(roomID), clientID.String()).Err()
}

func (p *RedisPresence) GetRoomMembers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error) {
	members, err := p.client.SMembers(ctx, roomClientsKey(roomID)).Result()
	if err != nil {
		return nil, err
	}
	clientIDs := make([]uuid.UUID, 0, len(members))
	for _, m := range members {
		id, err := uuid.Parse(m)
		if err != nil {
			continue
		}
		// Check if client is still connected
		connected, err := p.IsClientConnected(ctx, id)
		if err != nil || !connected {
			// Remove disconnected client from room
			p.client.SRem(ctx, roomClientsKey(roomID), m)
			continue
		}
		clientIDs = append(clientIDs, id)
	}
	return clientIDs, nil
}

func (p *RedisPresence) CleanupDisconnectedClients(ctx context.Context) error {
	// // This should be called periodically to clean up old client data
	// pattern := "client:*:status"
	// keys, err := p.client.Keys(ctx, pattern).Result()
	// if err != nil {
	// 	return err
	// }

	// for _, key := range keys {
	// 	data, err := p.client.Get(ctx, key).Result()
	// 	if err != nil {
	// 		continue
	// 	}

	// 	var status PresenceClient
	// 	if err := json.Unmarshal([]byte(data), &status); err != nil {
	// 		continue
	// 	}

	// 	// If client hasn't been seen for more than 1 hour, clean up
	// 	if time.Now().Unix()-status.LastSeen > 3600 {
	// 		clientID := status.ID

	// 		// Remove from any rooms
	// 		if status.RoomID != uuid.Nil {
	// 			p.LeaveRoom(ctx, status.RoomID, clientID)
	// 		}

	// 		// Remove client status and room mapping
	// 		p.client.Del(ctx, key)
	// 	}
	// }

	return nil
}

func (p *RedisPresence) UpdateHeartbeat(ctx context.Context, clientID uuid.UUID) error {
	status, err := p.GetClient(ctx, clientID)
	if err != nil {
		return err
	}

	status.LastSeen = time.Now().Unix()
	return p.SetClient(ctx, clientID, status)
}
