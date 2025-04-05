package main

import (
	"log"
	"net/http"

	"github.com/iulian509/realtime-messaging/config"
	"github.com/iulian509/realtime-messaging/publisher/internal/handlers"
	"github.com/iulian509/realtime-messaging/publisher/internal/mq"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load YAML configuration file: %v", err)
	}

	publisherClient, err := mq.NewPublisher(config.Nats.Host)
	if err != nil {
		log.Fatalf("failed to connect to NATS server: %v", err)
	}

	deps := &handlers.Dependencies{
		PublisherClient: publisherClient,
	}

	http.HandleFunc("/publish", deps.PublisherHandler)
	log.Println("publisher service running on :3000")
	err = http.ListenAndServe(":3000", nil)
	log.Fatalf("failed to start publisher service: %v", err)
}
