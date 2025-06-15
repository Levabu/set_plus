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
