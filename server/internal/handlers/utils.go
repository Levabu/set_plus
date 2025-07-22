package handlers

import (
	"fmt"
	"server/internal/domain"
)

func SendJSON(client *domain.LocalClient, payload interface{}) error {
	err := client.Conn.WriteJSON(payload)
	if err != nil {
		return fmt.Errorf("error sending message")
	}
	return nil
}

func SendError(client *domain.LocalClient, msg domain.ErrorMessage) error {
	return SendJSON(client, struct {
		Type domain.OutMessageType `json:"type"`
		domain.ErrorMessage
	}{
		Type:         domain.ErrorOut,
		ErrorMessage: msg,
	})
}
