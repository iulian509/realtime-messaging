package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iulian509/realtime-messaging/publisher/internal/mq"
)

const (
	pingPeriod = 30 * time.Second
	pongWait   = 60 * time.Second
	writeWait  = 30 * time.Second
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

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go pingConnection(ctx, cancel, conn)

	processMessages(ctx, cancel, conn, deps.PublisherClient)
}

func pingConnection(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait))
			if err != nil {
				log.Println("Ping error:", err)
				cancel()
				conn.Close()
				return
			}
		case <-ctx.Done():
			return
		}
	}
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
