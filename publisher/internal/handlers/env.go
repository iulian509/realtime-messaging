package handlers

import "github.com/iulian509/realtime-messaging/publisher/internal/mq"

type Env struct {
	PublisherClient *mq.Publisher
}
