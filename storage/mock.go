package storage

import "sync"

// MockStore is a mock implementation of StateStore for testing
type MockStore struct {
	mu              sync.Mutex
	ledStripCalls   []MockLEDStripCall
	ledBarCalls     []MockLEDBarCall
	videoLightCalls []MockVideoLightCall
}

// MockLEDStripCall represents a recorded LED strip save call
type MockLEDStripCall struct {
	ID int
	R  int
	G  int
	B  int
}

// MockLEDBarCall represents a recorded LED bar save call
type MockLEDBarCall struct {
	ID       int
	Channels []int
}

// MockVideoLightCall represents a recorded video light save call
type MockVideoLightCall struct {
	ID         int
	On         bool
	Brightness int
}

// NewMockStore creates a new mock store for testing
func NewMockStore() *MockStore {
	return &MockStore{
		ledStripCalls:   make([]MockLEDStripCall, 0),
		ledBarCalls:     make([]MockLEDBarCall, 0),
		videoLightCalls: make([]MockVideoLightCall, 0),
	}
}

// SaveLEDStripState records an LED strip save call
func (m *MockStore) SaveLEDStripState(id int, r, g, b int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ledStripCalls = append(m.ledStripCalls, MockLEDStripCall{
		ID: id,
		R:  r,
		G:  g,
		B:  b,
	})
	return nil
}

// LoadLEDStripState returns default values (mock doesn't persist)
func (m *MockStore) LoadLEDStripState(id int) (r, g, b int, err error) {
	// Mock always returns defaults
	return 0, 0, 0, nil
}

// SaveLEDBarChannels records an LED bar save call
func (m *MockStore) SaveLEDBarChannels(ledbarID int, channels []int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Make a copy of channels to avoid aliasing issues
	channelsCopy := make([]int, len(channels))
	copy(channelsCopy, channels)

	m.ledBarCalls = append(m.ledBarCalls, MockLEDBarCall{
		ID:       ledbarID,
		Channels: channelsCopy,
	})
	return nil
}

// LoadLEDBarChannels returns default values (mock doesn't persist)
func (m *MockStore) LoadLEDBarChannels(ledbarID int) ([]int, error) {
	// Mock always returns 77 zeros
	return make([]int, 77), nil
}

// SaveVideoLightState records a video light save call
func (m *MockStore) SaveVideoLightState(id int, on bool, brightness int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.videoLightCalls = append(m.videoLightCalls, MockVideoLightCall{
		ID:         id,
		On:         on,
		Brightness: brightness,
	})
	return nil
}

// LoadVideoLightState returns default values (mock doesn't persist)
func (m *MockStore) LoadVideoLightState(id int) (on bool, brightness int, err error) {
	// Mock always returns defaults
	return false, 0, nil
}

// GetLEDStripCalls returns all recorded LED strip save calls
func (m *MockStore) GetLEDStripCalls() []MockLEDStripCall {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]MockLEDStripCall, len(m.ledStripCalls))
	copy(result, m.ledStripCalls)
	return result
}

// GetLEDBarCalls returns all recorded LED bar save calls
func (m *MockStore) GetLEDBarCalls() []MockLEDBarCall {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]MockLEDBarCall, len(m.ledBarCalls))
	copy(result, m.ledBarCalls)
	return result
}

// GetVideoLightCalls returns all recorded video light save calls
func (m *MockStore) GetVideoLightCalls() []MockVideoLightCall {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make([]MockVideoLightCall, len(m.videoLightCalls))
	copy(result, m.videoLightCalls)
	return result
}

// Clear clears all recorded calls
func (m *MockStore) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ledStripCalls = make([]MockLEDStripCall, 0)
	m.ledBarCalls = make([]MockLEDBarCall, 0)
	m.videoLightCalls = make([]MockVideoLightCall, 0)
}
