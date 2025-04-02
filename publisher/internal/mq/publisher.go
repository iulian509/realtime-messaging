package mq

import (
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	conn *nats.Conn
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
