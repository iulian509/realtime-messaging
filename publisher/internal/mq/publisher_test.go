package mq

import (
	"testing"

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

func TestNewPublisherInvalidConnection(t *testing.T) {
	_, err := NewPublisher("nats://test:1234")
	assert.Error(t, err)
}

func TestPublishMessage(t *testing.T) {
	mockedConn := new(MockNATSConnection)
	publisher := &Publisher{Conn: mockedConn}

	const (
		publish     = "Publish"
		subject     = "subject"
		testMessage = "test message"
	)

	t.Run("successful publish", func(t *testing.T) {
		mockedConn.On(publish, subject, []byte(testMessage)).Return(nil).Once()

		err := publisher.PublishMessage(subject, []byte(testMessage))
		assert.NoError(t, err)
		mockedConn.AssertExpectations(t)
	})
	t.Run("failing publish", func(t *testing.T) {
		mockedConn.On(publish, subject, []byte(testMessage)).Return(assert.AnError).Once()

		err := publisher.PublishMessage(subject, []byte(testMessage))
		assert.Error(t, err)
		mockedConn.AssertExpectations(t)
	})
}
