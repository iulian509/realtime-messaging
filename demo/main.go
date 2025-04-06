package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	publisherURL    = "ws://publisher:3000/publish"
	subscribeURL    = "ws://subscriber:3001/subscribe"
	messageInterval = 1 * time.Second
	maxMessages     = 10
)

func main() {
	jwtToken := os.Getenv("JWT_TOKEN")
	if jwtToken == "" {
		log.Fatalf("JWT_TOKEN environment variable not set")
	}

	var wg sync.WaitGroup

	// wait for subscriber to connect and to be ready
	wg.Add(1)
	go startSubscriber(jwtToken, &wg)
	wg.Wait()

	publishMessages(jwtToken, &wg)

	// wait for all messages to be received
	wg.Wait()

	log.Println("Demo completed: all 10 messages sent and received")
}

func connectWebSocket(url, jwtToken string) (*websocket.Conn, error) {
	dialer := websocket.Dialer{}
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

	conn, _, err := dialer.Dial(url, header)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket endpoint %s: %v", url, err)
	}
	return conn, nil
}

func startSubscriber(jwtToken string, wg *sync.WaitGroup) {
	conn, err := connectWebSocket(subscribeURL, jwtToken)
	if err != nil {
		log.Fatalf("failed to connect to subscribe endpoint: %v", err)
	}
	defer conn.Close()

	log.Println("connected to subscribe ws, waiting for messages...")
	wg.Done()

	count := 0
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("ws read error: %v", err)
			break
		}
		fmt.Printf("received message from subscribe: %s\n", string(message))
		count++
		if count == maxMessages {
			wg.Done()
			break
		}
	}
}

func publishMessages(jwtToken string, wg *sync.WaitGroup) {
	conn, err := connectWebSocket(publisherURL, jwtToken)
	if err != nil {
		log.Fatalf("failed to connect to publish ws: %v", err)
	}
	defer conn.Close()

	wg.Add(1)
	for i := range [maxMessages]int{} {
		message := fmt.Sprintf("demo message #%d", i)
		fmt.Printf("sending message: %s\n", message)

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("failed to send message via publish WebSocket: %v", err)
		} else {
			log.Println("message published successfully")
		}

		time.Sleep(messageInterval)
	}
}
