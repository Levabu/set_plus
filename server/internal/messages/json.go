package messages

import (
	// "encoding/json"
	"fmt"
	"server/internal/server"
)

func sendJSON(client *server.Client, payload interface{}) error {
	// data, err := json.Marshal(payload)
	// if err != nil {
	// 	return fmt.Errorf("error marshalling json: %s", err.Error())
	// }

	err := client.Conn.WriteJSON(payload)
	if err != nil {
		return fmt.Errorf("error sending message")
	}
	return nil
}
