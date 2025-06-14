package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID uuid.UUID
	Conn *websocket.Conn
}