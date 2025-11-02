package presence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryPresence struct {
	clients map[uuid.UUID]PresenceClient
	roomClients map[uuid.UUID]map[uuid.UUID]struct{}
	mu sync.RWMutex
}

func NewMemoryPresence() *MemoryPresence {
	return &MemoryPresence{
		clients: make(map[uuid.UUID]PresenceClient),
		roomClients: make(map[uuid.UUID]map[uuid.UUID]struct{}),
	}
}

func (p *MemoryPresence) GetClient(ctx context.Context, clientID uuid.UUID) (PresenceClient, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	client, ok := p.clients[clientID]
	if ok {
		return client, nil
	}
	return PresenceClient{}, fmt.Errorf("client %s not found", clientID)
}

func (p *MemoryPresence) SetClient(ctx context.Context, clientID uuid.UUID, status PresenceClient) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	status.LastSeen = time.Now().Unix()
	p.clients[clientID] = status
	return nil
}
func (p *MemoryPresence) GetActiveRoomMembers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	clientIDs := make([]uuid.UUID, 0)
	room := p.roomClients[roomID]
	for clientID, _ := range room {
		clientIDs = append(clientIDs, clientID)
	}
	return clientIDs, nil
}

func (p *MemoryPresence) JoinRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error {
	if err := p.SetClient(ctx, clientID, PresenceClient{
		ID:        clientID,
		RoomID:    roomID,
		Connected: true,
		LastSeen:  time.Now().Unix(),
	}); err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.roomClients[roomID]; !ok {
		p.roomClients[roomID] = make(map[uuid.UUID]struct{})
	}
	p.roomClients[roomID][clientID] = struct{}{}
	return nil
}

func (p *MemoryPresence) LeaveRoom(ctx context.Context, clientID uuid.UUID) error {
	client, err := p.GetClient(ctx, clientID)
	if err != nil {
		return err
	}
	client.Connected = false
	client.LastSeen = time.Now().Unix()
	p.SetClient(ctx, clientID, client)
	roomClients := p.roomClients[client.RoomID]
	delete(roomClients, clientID)
	return nil
}

func (p *MemoryPresence) CleanupPresenceRoom(ctx context.Context, roomID uuid.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.roomClients, roomID)
	for id, client := range p.clients {
		if client.RoomID == roomID {
			delete(p.clients, id)
		}
	}
}