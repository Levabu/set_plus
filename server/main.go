package main

import (
	"encoding/json"
	"log"
	"net/http"
	"server/internal/messages"
	"server/internal/server"

	// "github.com/google/uuid"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // allow all origins (dev only)
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &server.Client{
		ID: uuid.New(),
		Conn: conn,
	}

	go handleConnetion(client)
}

func handleConnetion(client *server.Client) {
	defer client.Conn.Close()
	// fmt.Println(conn)

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		log.Printf("Received: %s", msg)

		var baseMessage messages.BaseMessage
		if err := json.Unmarshal(msg, &baseMessage); err != nil {
			log.Println("invalid message:", err)
			continue
		}

		log.Println("message type:", baseMessage.Type)
		handler, ok := messages.MessageHandlers[baseMessage.Type]
		if !ok {
			log.Println("unknown message type:", baseMessage.Type)
			continue
		}

		err = handler(client, msg)
		if err != nil {
			log.Println("error handling message:", err.Error())
		}

		// err = client.Conn.WriteMessage(websocket.TextMessage, []byte("Echo: "+string(msg)))
		// if err != nil {
		// 	log.Println("Write error:", err)
		// 	break
		// }
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// game, err := game.NewGame(game.Classic)
	// if err != nil {
	// 	log.Fatalf("Error creating game: %v", err)
	// }
	// game.GenerateCards()

	// for range 10 {
	// 	set := game.FindSet()
	// 	fmt.Println(set)
	// 	for _, card := range set {
	// 		fmt.Printf("Features: %v\n", card)
	// 	}
	// 	ids := make([]uuid.UUID, 0)
	// 	for _, card := range set {
	// 		ids = append(ids, card.CardID)
	// 	}
	// 	game.HandleCheckSet(ids)
	// }

	// game.HandleCheckSet(ids)

	// for _, card := range game.Deck {
	// 	fmt.Printf("Features: %v\n", card.IsDiscarded)
	// }

	// fmt.Println()
	// set = game.FindSet()
	// for _, card := range set {
	// 	fmt.Printf("Features: %v\n", card.Features)
	// }

	// log.Println("Game initialized with", len(game.Deck), "cards.")
	// log.Println("Game version:", game.GameConfig.Features, "Variations:", game.GameConfig.VariationsNumber)
}
