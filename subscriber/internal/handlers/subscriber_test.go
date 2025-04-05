package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iulian509/realtime-messaging/subscriber/internal/mq"
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

func setupTest(t *testing.T, mockConn *MockNATSConnection) (*httptest.Server, *websocket.Conn) {
	subscriber := &mq.Subscriber{Conn: mockConn}
	deps := &Dependencies{SubscriberClient: subscriber}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deps.SubscriberHandler(w, r)
	}))
	t.Cleanup(server.Close)

	url := "ws" + server.URL[4:]
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	t.Cleanup(func() { conn.Close() })

	return server, conn
}

func TestSubscriberHandler(t *testing.T) {
	t.Run("successful message receive", func(t *testing.T) {
		// Set up mocks
		mockedConn := new(MockNATSConnection)
		subscription := &nats.Subscription{}
		testMessage := []byte("test message")

		mockedConn.On("Subscribe", "subject", mock.Anything).Run(func(args mock.Arguments) {
			handler := args.Get(1).(nats.MsgHandler)
			go handler(&nats.Msg{Data: testMessage, Subject: "subject"})
		}).Return(subscription, nil)

		mockedConn.On("Flush").Return(nil)
		mockedConn.On("LastError").Return(nil)

		// Set up test server
		_, conn := setupTest(t, mockedConn)

		// Wait for message processing
		time.Sleep(100 * time.Millisecond)

		_, message, err := conn.ReadMessage()
		assert.NoError(t, err)
		assert.Equal(t, testMessage, message)

		mockedConn.AssertExpectations(t)
	})
	t.Run("subscription failure", func(t *testing.T) {
		// Set up mocks
		mockedConn := new(MockNATSConnection)
		mockedConn.On("Subscribe", "subject", mock.Anything).Return(nil, assert.AnError)

		// Set up test server
		_, conn := setupTest(t, mockedConn)

		// Wait for the subscription failure to be handled
		time.Sleep(100 * time.Millisecond)

		_, _, err := conn.ReadMessage()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "close")

		mockedConn.AssertExpectations(t)
	})
}
