package mq

import (
	"github.com/nats-io/nats.go"
)

type NATSConnection interface {
	Publish(subj string, data []byte) error
}

type Publisher struct {
	Conn NATSConnection
}

func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Publisher{Conn: nc}, nil
}

func (p *Publisher) PublishMessage(message []byte) error {
	err := p.Conn.Publish("subject", message)
	if err != nil {
		return err
	}

	return nil
}
