package handlers

import "github.com/iulian509/realtime-messaging/subscriber/internal/mq"

type Dependencies struct {
	SubscriberClient *mq.Subscriber
}
