package domain

import (
	"encoding/json"
)

type MessageRouter interface {
	RegisterHandler(msgType InMessageType, handler MessageHandler)
	HandleMessage(client *LocalClient, msgType InMessageType, rawMsg json.RawMessage) error
}

type ConnectionManager interface {
	HandleConnection(client *LocalClient)
}
