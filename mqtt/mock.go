package mqtt

import "sync"

// MockPublisher is a mock implementation of the Publisher interface for testing
type MockPublisher struct {
	mu       sync.Mutex
	messages []Message
}

// Message represents a published MQTT message
type Message struct {
	Topic   string
	Payload interface{}
}

// NewMockPublisher creates a new mock publisher for testing
func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		messages: make([]Message, 0),
	}
}

// Publish records the published message
func (m *MockPublisher) Publish(topic string, payload interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, Message{
		Topic:   topic,
		Payload: payload,
	})

	return nil
}

// GetMessages returns all published messages
func (m *MockPublisher) GetMessages() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a copy to avoid race conditions
	result := make([]Message, len(m.messages))
	copy(result, m.messages)
	return result
}

// GetLastMessage returns the last published message
func (m *MockPublisher) GetLastMessage() *Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.messages) == 0 {
		return nil
	}

	return &m.messages[len(m.messages)-1]
}

// Clear clears all recorded messages
func (m *MockPublisher) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = make([]Message, 0)
}

// MessageCount returns the number of messages published
func (m *MockPublisher) MessageCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.messages)
}
