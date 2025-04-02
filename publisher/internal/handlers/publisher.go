package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func (env *Env) PublisherHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("client connected to WebSocket")

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err) {
				log.Printf("WebSocket connection closed unexpectedly: %v", err)
			} else {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}
		log.Printf("received message: %s", message)

		err = env.PublisherClient.PublishMessage(message)
		if err != nil {
			log.Printf("error publishing message: %v", err)
		} else {
			log.Printf("published message: %s", message)
		}
	}
}
