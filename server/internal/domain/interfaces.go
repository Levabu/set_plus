package domain

import (
	"encoding/json"
)

type MessageRouter interface {
	RegisterHandler(msgType InMessageType, handler MessageHandler)
	HandleMessage(client *Client, msgType InMessageType, rawMsg json.RawMessage) error
}

type ConnectionManager interface {
	HandleConnection(client *Client)
}