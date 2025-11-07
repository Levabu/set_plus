package domain

import (
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type LocalClient struct {
	ID             uuid.UUID
	Conn           *websocket.Conn
	WriteChan      chan interface{}
	RoomID         uuid.UUID
	Nickname       string
	Connected      bool
	DisconnectedAt time.Time
	ReconnectTimer *time.Timer
}

type LocalClientManager interface {
	Add(client *LocalClient)
	Get(id uuid.UUID) *LocalClient
	Remove(id uuid.UUID)
	GetAll() map[uuid.UUID]*LocalClient
	SetClientConnected(id uuid.UUID, connected bool)
	CleanupLocalRoomClients(roomID uuid.UUID)
	IsRoomEmpty(roomID uuid.UUID) bool
}

type LocalClients struct {
	clients map[uuid.UUID]*LocalClient
	mu      sync.RWMutex
}

func NewLocalClients() *LocalClients {
	return &LocalClients{
		clients: make(map[uuid.UUID]*LocalClient),
	}
}

func (c *LocalClients) Add(client *LocalClient) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[client.ID] = client
}

func (c *LocalClients) Get(id uuid.UUID) *LocalClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.clients[id]
}

func (c *LocalClients) Remove(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, id)
}

func (c *LocalClients) SetClientConnected(id uuid.UUID, connected bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	client, ok := c.clients[id]
	if !ok {
		return
	}
	client.Connected = connected
	if !connected {
		client.DisconnectedAt = time.Now()
	}
	c.clients[id] = client
}

func (c *LocalClients) GetAll() map[uuid.UUID]*LocalClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[uuid.UUID]*LocalClient)
	for id, client := range c.clients {
		result[id] = client
	}
	return result
}

func (c *LocalClients) CleanupLocalRoomClients(roomID uuid.UUID) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for id, client := range c.clients {
		if client.RoomID == roomID {
			delete(c.clients, id)
		}
	}
}

func (c *LocalClients) IsRoomEmpty(roomID uuid.UUID) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	counter := 0
	for _, client := range c.clients {
		if client.RoomID == roomID {
			counter++
		}
	}

	return counter == 0
}
