package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iulian509/realtime-messaging/publisher/internal/mq"

	iw "github.com/iulian509/realtime-messaging/internal/websocket"
)

var upgrader = websocket.Upgrader{}

func (deps *Dependencies) PublisherHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("client connected to Publisher WebSocket")

	iw.SetHeartbeatConfig(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go iw.PingConnection(ctx, cancel, conn)

	processMessages(ctx, cancel, conn, deps.PublisherClient)
}

func processMessages(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, publisherClient *mq.Publisher) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err) {
					log.Printf("WebSocket connection closed unexpectedly: %v", err)
				} else {
					log.Printf("WebSocket read error: %v", err)
				}
				cancel()
				return
			}
			log.Printf("received message: %s", message)

			err = publisherClient.PublishMessage(message)
			if err != nil {
				log.Printf("error publishing message: %v", err)
			} else {
				log.Printf("published message: %s", message)
			}
		}
	}
}
