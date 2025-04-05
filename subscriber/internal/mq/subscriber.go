package mq

import "github.com/nats-io/nats.go"

type NATSConnection interface {
	Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error)
	Flush() error
	LastError() error
	Close()
}

type NATSMessageHandler func(msg *nats.Msg)

type Subscriber struct {
	Conn NATSConnection
}

func NewSubscriber(url string) (*Subscriber, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Subscriber{Conn: nc}, nil
}

func (s *Subscriber) Subscribe(subject string, handler NATSMessageHandler) (*nats.Subscription, error) {
	sub, err := s.Conn.Subscribe(subject, nats.MsgHandler(handler))
	if err != nil {
		return nil, err
	}

	if err := s.Conn.Flush(); err != nil {
		return nil, err
	}
	if err := s.Conn.LastError(); err != nil {
		return nil, err
	}

	return sub, nil
}
