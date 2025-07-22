package domain

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type LocalClient struct {
	ID       uuid.UUID
	Conn     *websocket.Conn
	RoomID   uuid.UUID
	Nickname string
}

type LocalClientManager interface {
	Add(client *LocalClient)
	Get(id uuid.UUID) *LocalClient
	Remove(id uuid.UUID)
	GetAll() map[uuid.UUID]*LocalClient
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

func (c *LocalClients) GetAll() map[uuid.UUID]*LocalClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[uuid.UUID]*LocalClient)
	for id, client := range c.clients {
		result[id] = client
	}
	return result
}
