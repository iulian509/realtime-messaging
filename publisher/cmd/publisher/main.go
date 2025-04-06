package main

import (
	"log"
	"net/http"

	"github.com/iulian509/realtime-messaging/config"
	"github.com/iulian509/realtime-messaging/internal/auth"
	"github.com/iulian509/realtime-messaging/internal/metrics"
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

	metrics.InitMetrics()

	deps := &handlers.Dependencies{
		PublisherClient: publisherClient,
	}

	http.HandleFunc("/metrics", metrics.PromHandler())
	http.HandleFunc("/publish", metrics.PrometheusMiddleware(auth.AuthMiddleware(deps.PublisherHandler)))

	err = http.ListenAndServe(":3000", nil)
	log.Println("publisher service running on :3000")
	if err != nil {
		log.Fatalf("failed to start publisher service: %v", err)
	}
}
