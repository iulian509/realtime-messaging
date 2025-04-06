package main

import (
	"log"
	"net/http"

	"github.com/iulian509/realtime-messaging/internal/auth"
	"github.com/iulian509/realtime-messaging/internal/config"
	"github.com/iulian509/realtime-messaging/internal/metrics"
	"github.com/iulian509/realtime-messaging/subscriber/internal/handlers"
	"github.com/iulian509/realtime-messaging/subscriber/internal/mq"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load YAML configuration file: %v", err)
	}

	subscriberClient, err := mq.NewSubscriber(config.Nats.Host)
	if err != nil {
		log.Fatalf("failed to connect to NATS server: %v", err)
	}

	metrics.InitMetrics()

	deps := &handlers.Dependencies{
		SubscriberClient: subscriberClient,
	}

	http.HandleFunc("/metrics", metrics.PromHandler())
	http.HandleFunc("/subscribe", metrics.PrometheusMiddleware(auth.AuthMiddleware(deps.SubscriberHandler)))

	err = http.ListenAndServe(":3001", nil)
	log.Println("subscriber service running on :3001")
	if err != nil {
		log.Fatalf("failed to start subscriber service: %v", err)
	}
}
