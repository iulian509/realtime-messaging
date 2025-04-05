package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/iulian509/realtime-messaging/publisher/internal/mq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockNATSConnection struct {
	mock.Mock
}

func (m *MockNATSConnection) Publish(subj string, data []byte) error {
	args := m.Called(subj, data)
	return args.Error(0)
}

func TestPublisherHandler(t *testing.T) {
	// Set up mocks
	mockedConn := new(MockNATSConnection)
	publisher := &mq.Publisher{Conn: mockedConn}
	deps := &Dependencies{PublisherClient: publisher}

	// Set up test handler and server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		deps.PublisherHandler(w, r)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	// Upgrade connection to WebSocket
	url := "ws" + server.URL[4:]
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)
	assert.NoError(t, err)
	defer ws.Close()

	const (
		publish     = "Publish"
		subject     = "subject"
		testMessage = "test message"
	)

	t.Run("successfull message publishing", func(t *testing.T) {
		mockedConn.On(publish, subject, mock.AnythingOfType("[]uint8")).Return(nil).Once()

		err = ws.WriteMessage(websocket.TextMessage, []byte(testMessage))
		assert.NoError(t, err)

		// Wait for the handler to process the message
		time.Sleep(100 * time.Millisecond)

		mockedConn.AssertExpectations(t)
	})
	t.Run("message publishing failure", func(t *testing.T) {
		mockedConn.On(publish, subject, mock.AnythingOfType("[]uint8")).Return(assert.AnError).Once()

		err = ws.WriteMessage(websocket.TextMessage, []byte(testMessage))
		assert.NoError(t, err, "message writing should succeed even if NATS publish fails")

		// Wait for the handler to process the message
		time.Sleep(100 * time.Millisecond)

		mockedConn.AssertExpectations(t)
	})
}
