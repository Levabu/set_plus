package presence

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryPresence struct {
	clients           map[uuid.UUID]PresenceClient
	activeRoomClients map[uuid.UUID]map[uuid.UUID]struct{}
	mu                sync.RWMutex
}

func NewMemoryPresence() *MemoryPresence {
	return &MemoryPresence{
		clients:           make(map[uuid.UUID]PresenceClient),
		activeRoomClients: make(map[uuid.UUID]map[uuid.UUID]struct{}),
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

	status.DisconnectedAt = time.Now().Unix()
	p.clients[clientID] = status
	return nil
}
func (p *MemoryPresence) GetActiveRoomMembers(ctx context.Context, roomID uuid.UUID) ([]uuid.UUID, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	clientIDs := make([]uuid.UUID, 0)
	room := p.activeRoomClients[roomID]
	for clientID, _ := range room {
		clientIDs = append(clientIDs, clientID)
	}
	return clientIDs, nil
}

func (p *MemoryPresence) JoinRoom(ctx context.Context, roomID uuid.UUID, clientID uuid.UUID) error {
	if err := p.SetClient(ctx, clientID, PresenceClient{
		ID:             clientID,
		RoomID:         roomID,
		Connected:      true,
	}); err != nil {
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.activeRoomClients[roomID]; !ok {
		p.activeRoomClients[roomID] = make(map[uuid.UUID]struct{})
	}
	p.activeRoomClients[roomID][clientID] = struct{}{}
	return nil
}

func (p *MemoryPresence) LeaveRoom(ctx context.Context, clientID uuid.UUID) error {
	client, err := p.GetClient(ctx, clientID)
	if err != nil {
		return err
	}
	client.Connected = false
	client.DisconnectedAt = time.Now().Unix()
	p.SetClient(ctx, clientID, client)
	roomClients, ok := p.activeRoomClients[client.RoomID]
	if ok {
		delete(roomClients, clientID)
	}
	return nil
}

func (p *MemoryPresence) CleanupPresenceRoom(ctx context.Context, roomID uuid.UUID) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.activeRoomClients, roomID)
	for id, client := range p.clients {
		if client.RoomID == roomID {
			delete(p.clients, id)
		}
	}
}

func (p *MemoryPresence) RemoveClient(ctx context.Context, clientID uuid.UUID, roomID uuid.UUID) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if activeClients, ok := p.activeRoomClients[roomID]; ok {
		delete(activeClients, clientID)
	}
	delete(p.clients, clientID)
	return nil
}