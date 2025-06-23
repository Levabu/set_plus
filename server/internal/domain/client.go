package domain

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID       uuid.UUID
	Conn     *websocket.Conn
	RoomID   uuid.UUID
	Nickname string
}

type ClientManager interface {
	Add(client *Client)
	Get(id uuid.UUID) *Client
	Remove(id uuid.UUID)
	GetAll() map[uuid.UUID]*Client
}

type LocalClients struct {
	clients map[uuid.UUID]*Client
	mu      sync.RWMutex
}

func NewLocalClients() *LocalClients {
	return &LocalClients{
		clients: make(map[uuid.UUID]*Client),
	}
}

func (c *LocalClients) Add(client *Client) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[client.ID] = client
}

func (c *LocalClients) Get(id uuid.UUID) *Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.clients[id]
}

func (c *LocalClients) Remove(id uuid.UUID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, id)
}

func (c *LocalClients) GetAll() map[uuid.UUID]*Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[uuid.UUID]*Client)
	for id, client := range c.clients {
		result[id] = client
	}
	return result
}