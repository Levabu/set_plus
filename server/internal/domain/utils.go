package domain

import (
	"errors"
)

func SendJSON(client *LocalClient, payload interface{}) error {
	select {
	case client.WriteChan <- payload:
		return nil
	default:
		return errors.New("write channel full")
	}
	
}

func SendError(client *LocalClient, msg ErrorMessage) error {
	return SendJSON(client, struct {
		Type OutMessageType `json:"type"`
		ErrorMessage
	}{
		Type:         ErrorOut,
		ErrorMessage: msg,
	})
}
