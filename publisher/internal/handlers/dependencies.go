package handlers

import "github.com/iulian509/realtime-messaging/publisher/internal/mq"

type Dependencies struct {
	PublisherClient *mq.Publisher
}
