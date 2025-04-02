package mq

import (
	"github.com/nats-io/nats.go"
)

type NATSConnection interface {
	Publish(subj string, data []byte) error
}

type Publisher struct {
	conn NATSConnection
}

func NewPublisher(url string) (*Publisher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Publisher{conn: nc}, nil
}

func (p *Publisher) PublishMessage(message []byte) error{
	err := p.conn.Publish("subject", message)
	if err != nil {
		return err
	}

	return nil
}
