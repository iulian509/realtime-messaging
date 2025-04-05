package mq

import (
	"errors"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNATSConnection struct {
	mock.Mock
}

func (m *MockNATSConnection) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	args := m.Called(subj, cb)
	return args.Get(0).(*nats.Subscription), args.Error(1)
}

func (m *MockNATSConnection) Flush() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNATSConnection) LastError() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockNATSConnection) Close() {
	m.Called()
}

func TestNewSubscriberInvalidConnection(t *testing.T) {
	_, err := NewSubscriber("nats://test:1234")
	assert.Error(t, err)
}

func TestSubscriber(t *testing.T) {
	const (
		subscribe = "Subscribe"
		subject   = "subject"
	)

	dummySubscription := &nats.Subscription{}

	setupMockedConn := func(subscribeError, flushError, lastError error) *MockNATSConnection {
		mockedConn := new(MockNATSConnection)

		if subscribeError != nil {
			mockedConn.On(subscribe, mock.Anything, mock.Anything).Return((*nats.Subscription)(nil), subscribeError)
		} else {
			mockedConn.On(subscribe, mock.Anything, mock.Anything).Return(dummySubscription, nil)
			mockedConn.On("Flush").Return(flushError)
			mockedConn.On("LastError").Maybe().Return(lastError)
		}

		return mockedConn
	}

	t.Run("successful subscription", func(t *testing.T) {
		mockedConn := setupMockedConn(nil, nil, nil)
		subscriber := &Subscriber{Conn: mockedConn}

		sub, err := subscriber.Subscribe(subject, func(msg *nats.Msg) {})

		assert.NoError(t, err)
		assert.NotNil(t, sub)
		mockedConn.AssertExpectations(t)
	})

	t.Run("subscribe failure", func(t *testing.T) {
		mockConn := setupMockedConn(errors.New("subscription failed"), nil, nil)
		subscriber := &Subscriber{Conn: mockConn}

		sub, err := subscriber.Subscribe(subject, func(msg *nats.Msg) {})

		assert.Error(t, err)
		assert.Equal(t, "subscription failed", err.Error())
		assert.Nil(t, sub)
		mockConn.AssertExpectations(t)
	})

	t.Run("flush failure", func(t *testing.T) {
		mockConn := setupMockedConn(nil, errors.New("flush failed"), nil)
		subscriber := &Subscriber{Conn: mockConn}

		sub, err := subscriber.Subscribe(subject, func(msg *nats.Msg) {})

		assert.Error(t, err)
		assert.Equal(t, "flush failed", err.Error())
		assert.Nil(t, sub)
		mockConn.AssertExpectations(t)
	})

	t.Run("last error failure", func(t *testing.T) {
		mockConn := setupMockedConn(nil, nil, errors.New("last error occurred"))
		subscriber := &Subscriber{Conn: mockConn}

		sub, err := subscriber.Subscribe(subject, func(msg *nats.Msg) {})

		assert.Error(t, err)
		assert.Equal(t, "last error occurred", err.Error())
		assert.Nil(t, sub)
		mockConn.AssertExpectations(t)
	})
}
