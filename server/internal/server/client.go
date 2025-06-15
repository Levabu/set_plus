package server

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

type LocalClients struct {
	Clients map[uuid.UUID]*Client
	Mu      sync.RWMutex
}

func NewLocalClients() *LocalClients {
	return &LocalClients{
		Clients: make(map[uuid.UUID]*Client),
	}
}

func (c *LocalClients) Add(client *Client) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Clients[client.ID] = client
}

func (c *LocalClients) Get(id uuid.UUID) *Client {
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	client, ok := c.Clients[id]
	if !ok {
		return nil
	}
	return client
}

func (c *LocalClients) Remove(id uuid.UUID) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	delete(c.Clients, id)
}