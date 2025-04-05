package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/iulian509/realtime-messaging/subscriber/internal/mq"
	"github.com/nats-io/nats.go"

	iw "github.com/iulian509/realtime-messaging/internal/websocket"
)

var upgrader = websocket.Upgrader{}

func (deps *Dependencies) SubscriberHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("client connected to Subscriber WebSocket")

	iw.SetHeartbeatConfig(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go iw.PingConnection(ctx, cancel, conn)

	processMessages(ctx, cancel, conn, deps.SubscriberClient)
}

func processMessages(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn, subscriberClient *mq.Subscriber) {
	const subject = "subject"

	subscription, err := subscriberClient.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("received message on [%s]: %s", msg.Subject, string(msg.Data))

		err := conn.WriteMessage(websocket.TextMessage, msg.Data)
		if err != nil {
			log.Printf("failed to send message to WebSocket: %v", err)
			cancel()
			return
		}
	})
	if err != nil {
		log.Printf("failed to subscribe to subject %q: %v", subject, err)
		cancel()
		return
	}

	defer func() {
		err := subscription.Unsubscribe()
		if err != nil {
			log.Printf("error unsubscribing from subject %q: %v", subject, err)
		} else {
			log.Printf("unsubscribed from subject %q", subject)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)
				cancel()
				return
			}
		}
	}
}
