package messages

import (
	"fmt"
	"server/internal/server"
)

func SendJSON(client *server.Client, payload interface{}) error {
	err := client.Conn.WriteJSON(payload)
	if err != nil {
		return fmt.Errorf("error sending message")
	}
	return nil
}

func SendError(client *server.Client, msg ErrorMessage) error {
	return SendJSON(client, struct{ 
		Type string `json:"error"`
		ErrorMessage
	}{
		Type: "error",
		ErrorMessage: msg,
	})
}