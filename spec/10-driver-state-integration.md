# Driver State Integration

## Overview
Modify existing light drivers to integrate with the state storage layer, enabling automatic state persistence.

## Design Pattern

### Storage Interface
Drivers should accept a storage interface rather than a concrete implementation. This allows:
- Easy testing with mock storage
- Flexibility to change storage backend
- Clean separation of concerns

```go
type StateStore interface {
    SaveLEDStripState(id int, r, g, b int) error
    SaveLEDBarChannels(ledbarID int, channels []int) error
    SaveVideoLightState(id int, on bool, brightness int) error
}
```

## LED Strip Driver Modifications

### File: `drivers/ledstrip/ledstrip.go`

#### Add Storage Fields
```go
type LEDStrip struct {
    r         int
    g         int
    b         int
    publisher Publisher
    topic     string
    store     StateStore  // Add this
    id        int         // Add this (always 0 for now)
}
```

#### Update Constructor
Create new constructor that accepts initial state:

```go
// NewLEDStripWithState creates LED strip with initial state from storage
func NewLEDStripWithState(publisher Publisher, topic string, store StateStore, id int, r, g, b int) *LEDStrip {
    return &LEDStrip{
        r:         r,
        g:         g,
        b:         b,
        publisher: publisher,
        topic:     topic,
        store:     store,
        id:        id,
    }
}
```

Keep existing constructor for backward compatibility:
```go
// NewLEDStrip creates LED strip with default state (all off)
func NewLEDStrip(publisher Publisher, topic string) *LEDStrip {
    return NewLEDStripWithState(publisher, topic, nil, 0, 0, 0, 0)
}
```

#### Modify Publish Method
Update `Publish()` to save state after successful MQTT publish:

```go
func (l *LEDStrip) Publish() error {
    payload, err := l.formatMessage()
    if err != nil {
        return fmt.Errorf("failed to format message: %w", err)
    }

    if err := l.publisher.Publish(l.topic, payload); err != nil {
        return fmt.Errorf("failed to publish: %w", err)
    }

    // Save state to storage after successful publish
    if l.store != nil {
        if err := l.store.SaveLEDStripState(l.id, l.r, l.g, l.b); err != nil {
            // Log error but don't fail the operation
            // State will be out of sync, but light was updated
            log.Printf("Warning: Failed to save LED strip state: %v", err)
        }
    }

    return nil
}
```

#### Testing Considerations
- Update tests to pass `nil` for store parameter (existing behavior)
- Add new tests that verify state is saved after publish
- Use mock store to capture save calls

## LED Bar Driver Modifications

### File: `drivers/ledbar/ledbar.go`

#### Add Storage Fields
```go
type LEDBar struct {
    rgbw1     [6][4]int
    white1    [13]int
    rgbw2     [6][4]int
    white2    [13]int
    barID     int
    publisher Publisher
    topic     string
    store     StateStore  // Add this
}
```

#### Update Constructor
```go
// NewLEDBarWithState creates LED bar with initial state from storage
func NewLEDBarWithState(barID int, publisher Publisher, topic string, store StateStore, channels []int) (*LEDBar, error) {
    if barID < 0 {
        return nil, fmt.Errorf("barID must be non-negative, got %d", barID)
    }

    bar := &LEDBar{
        barID:     barID,
        publisher: publisher,
        topic:     topic,
        store:     store,
    }

    // Load state from channels array
    if err := bar.loadFromChannels(channels); err != nil {
        return nil, fmt.Errorf("failed to load channels: %w", err)
    }

    return bar, nil
}
```

Keep existing constructor:
```go
// NewLEDBar creates LED bar with default state (all off)
func NewLEDBar(barID int, publisher Publisher, topic string) (*LEDBar, error) {
    channels := make([]int, 77)
    return NewLEDBarWithState(barID, publisher, topic, nil, channels)
}
```

#### Add Helper Method: Load From Channels
```go
// loadFromChannels populates LED states from 77-value channel array
func (l *LEDBar) loadFromChannels(channels []int) error {
    if len(channels) != 77 {
        return fmt.Errorf("expected 77 channels, got %d", len(channels))
    }

    idx := 0

    // Load first section RGBW (24 values)
    for i := 0; i < 6; i++ {
        for j := 0; j < 4; j++ {
            l.rgbw1[i][j] = channels[idx]
            idx++
        }
    }

    // Load first section white (13 values)
    for i := 0; i < 13; i++ {
        l.white1[i] = channels[idx]
        idx++
    }

    // Skip 3 ignored values
    idx += 3

    // Load second section RGBW (24 values)
    for i := 0; i < 6; i++ {
        for j := 0; j < 4; j++ {
            l.rgbw2[i][j] = channels[idx]
            idx++
        }
    }

    // Load second section white (13 values)
    for i := 0; i < 13; i++ {
        l.white2[i] = channels[idx]
        idx++
    }

    return nil
}
```

#### Add Helper Method: Get Channels
```go
// getChannels returns current state as 77-value array
func (l *LEDBar) getChannels() []int {
    channels := make([]int, 77)
    idx := 0

    // First section RGBW
    for i := 0; i < 6; i++ {
        for j := 0; j < 4; j++ {
            channels[idx] = l.rgbw1[i][j]
            idx++
        }
    }

    // First section white
    for i := 0; i < 13; i++ {
        channels[idx] = l.white1[i]
        idx++
    }

    // 3 ignored values (already 0 from make)
    idx += 3

    // Second section RGBW
    for i := 0; i < 6; i++ {
        for j := 0; j < 4; j++ {
            channels[idx] = l.rgbw2[i][j]
            idx++
        }
    }

    // Second section white
    for i := 0; i < 13; i++ {
        channels[idx] = l.white2[i]
        idx++
    }

    return channels
}
```

#### Modify Publish Method
```go
func (l *LEDBar) Publish() error {
    payload := l.formatMessage()

    if err := l.publisher.Publish(l.topic, payload); err != nil {
        return fmt.Errorf("failed to publish: %w", err)
    }

    // Save state to storage after successful publish
    if l.store != nil {
        channels := l.getChannels()
        if err := l.store.SaveLEDBarChannels(l.barID, channels); err != nil {
            log.Printf("Warning: Failed to save LED bar state: %v", err)
        }
    }

    return nil
}
```

#### Testing Considerations
- Test `loadFromChannels()` with valid and invalid inputs
- Test `getChannels()` returns correct 77-value array
- Verify state save is called after publish
- Test round-trip: channels → load → save → channels

## Video Light Driver Modifications

### File: `drivers/videolight/videolight.go`

#### Add Storage Fields
```go
type VideoLight struct {
    on         bool
    brightness int
    lightID    int
    publisher  Publisher
    topic      string
    store      StateStore  // Add this
}
```

#### Update Constructor
```go
// NewVideoLightWithState creates video light with initial state from storage
func NewVideoLightWithState(lightID int, publisher Publisher, topic string, store StateStore, on bool, brightness int) (*VideoLight, error) {
    if lightID < 1 {
        return nil, fmt.Errorf("lightID must be positive, got %d", lightID)
    }

    if err := validateBrightness(brightness); err != nil {
        // Don't fail on invalid stored state, just fix it
        brightness = 0
        on = false
    }

    return &VideoLight{
        on:         on,
        brightness: brightness,
        lightID:    lightID,
        publisher:  publisher,
        topic:      topic,
        store:      store,
    }, nil
}
```

Keep existing constructor:
```go
// NewVideoLight creates video light with default state (off)
func NewVideoLight(lightID int, publisher Publisher, topic string) (*VideoLight, error) {
    return NewVideoLightWithState(lightID, publisher, topic, nil, false, 0)
}
```

#### Modify Publish Method
```go
func (v *VideoLight) Publish() error {
    payload := v.formatMessage()

    if err := v.publisher.Publish(v.topic, payload); err != nil {
        return fmt.Errorf("failed to publish: %w", err)
    }

    // Save state to storage after successful publish
    if v.store != nil {
        if err := v.store.SaveVideoLightState(v.lightID, v.on, v.brightness); err != nil {
            log.Printf("Warning: Failed to save video light state: %v", err)
        }
    }

    return nil
}
```

#### Handle Video Light IDs
Note: Video lights use IDs 0 and 1 in the database, but IDs 1 and 2 in the driver.

**Solution:** Map database IDs to driver IDs in main.go:
- Database ID 0 → Driver ID 1 (light 1)
- Database ID 1 → Driver ID 2 (light 2)

For storage operations in the driver, use `lightID - 1` when saving:

```go
if v.store != nil {
    dbID := v.lightID - 1  // Convert driver ID to database ID
    if err := v.store.SaveVideoLightState(dbID, v.on, v.brightness); err != nil {
        log.Printf("Warning: Failed to save video light state: %v", err)
    }
}
```

**Alternative:** Change database schema to use IDs 1 and 2 to match drivers.

## Storage Interface Mock for Testing

### File: `storage/mock.go`

```go
package storage

import "sync"

type MockStore struct {
    mu                 sync.Mutex
    ledStripCalls      []LEDStripState
    ledBarCalls        []LEDBarState
    videoLightCalls    []VideoLightState
}

type LEDStripState struct {
    ID  int
    R   int
    G   int
    B   int
}

type LEDBarState struct {
    ID       int
    Channels []int
}

type VideoLightState struct {
    ID         int
    On         bool
    Brightness int
}

func NewMockStore() *MockStore {
    return &MockStore{
        ledStripCalls:   make([]LEDStripState, 0),
        ledBarCalls:     make([]LEDBarState, 0),
        videoLightCalls: make([]VideoLightState, 0),
    }
}

func (m *MockStore) SaveLEDStripState(id int, r, g, b int) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.ledStripCalls = append(m.ledStripCalls, LEDStripState{
        ID: id,
        R:  r,
        G:  g,
        B:  b,
    })
    return nil
}

func (m *MockStore) SaveLEDBarChannels(ledbarID int, channels []int) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Make a copy of channels
    channelsCopy := make([]int, len(channels))
    copy(channelsCopy, channels)

    m.ledBarCalls = append(m.ledBarCalls, LEDBarState{
        ID:       ledbarID,
        Channels: channelsCopy,
    })
    return nil
}

func (m *MockStore) SaveVideoLightState(id int, on bool, brightness int) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.videoLightCalls = append(m.videoLightCalls, VideoLightState{
        ID:         id,
        On:         on,
        Brightness: brightness,
    })
    return nil
}

func (m *MockStore) GetLEDStripCalls() []LEDStripState {
    m.mu.Lock()
    defer m.mu.Unlock()

    result := make([]LEDStripState, len(m.ledStripCalls))
    copy(result, m.ledStripCalls)
    return result
}

func (m *MockStore) GetLEDBarCalls() []LEDBarState {
    m.mu.Lock()
    defer m.mu.Unlock()

    result := make([]LEDBarState, len(m.ledBarCalls))
    copy(result, m.ledBarCalls)
    return result
}

func (m *MockStore) GetVideoLightCalls() []VideoLightState {
    m.mu.Lock()
    defer m.mu.Unlock()

    result := make([]VideoLightState, len(m.videoLightCalls))
    copy(result, m.videoLightCalls)
    return result
}

func (m *MockStore) Clear() {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.ledStripCalls = make([]LEDStripState, 0)
    m.ledBarCalls = make([]LEDBarState, 0)
    m.videoLightCalls = make([]VideoLightState, 0)
}
```

## Testing Strategy

### LED Strip Driver Tests
```go
func TestLEDStripStatePersistence(t *testing.T) {
    mockPub := mqtt.NewMockPublisher()
    mockStore := storage.NewMockStore()

    strip := ledstrip.NewLEDStripWithState(mockPub, "test", mockStore, 0, 100, 150, 200)

    // Change color
    strip.SetColor(50, 75, 100)

    // Verify state was saved
    calls := mockStore.GetLEDStripCalls()
    if len(calls) != 1 {
        t.Errorf("Expected 1 save call, got %d", len(calls))
    }

    if calls[0].R != 50 || calls[0].G != 75 || calls[0].B != 100 {
        t.Errorf("Wrong values saved: %+v", calls[0])
    }
}
```

### LED Bar Driver Tests
```go
func TestLEDBarStatePersistence(t *testing.T) {
    mockPub := mqtt.NewMockPublisher()
    mockStore := storage.NewMockStore()

    // Create with initial channels
    initialChannels := make([]int, 77)
    initialChannels[0] = 100  // Set first channel

    bar, _ := ledbar.NewLEDBarWithState(0, mockPub, "test", mockStore, initialChannels)

    // Verify initial state loaded
    r, g, b, w, _ := bar.GetRGBW(1, 0)
    if r != 100 {
        t.Errorf("Initial state not loaded correctly")
    }

    // Change state
    bar.SetRGBW(1, 0, 255, 0, 0, 0)

    // Verify saved
    calls := mockStore.GetLEDBarCalls()
    if len(calls) != 1 {
        t.Errorf("Expected 1 save call, got %d", len(calls))
    }

    if calls[0].Channels[0] != 255 {
        t.Errorf("Wrong channel value saved")
    }
}
```

### Video Light Driver Tests
```go
func TestVideoLightStatePersistence(t *testing.T) {
    mockPub := mqtt.NewMockPublisher()
    mockStore := storage.NewMockStore()

    light, _ := videolight.NewVideoLightWithState(1, mockPub, "test", mockStore, true, 50)

    // Verify initial state
    on, brightness := light.GetState()
    if !on || brightness != 50 {
        t.Error("Initial state not loaded correctly")
    }

    // Change state
    light.TurnOn(75)

    // Verify saved (should save with dbID = 0, since lightID = 1)
    calls := mockStore.GetVideoLightCalls()
    if len(calls) != 1 {
        t.Errorf("Expected 1 save call, got %d", len(calls))
    }

    if calls[0].Brightness != 75 || !calls[0].On {
        t.Errorf("Wrong values saved: %+v", calls[0])
    }
}
```

## Success Criteria
- ✅ All drivers support state storage via StateStore interface
- ✅ State is saved after every successful MQTT publish
- ✅ Initial state can be loaded via new constructors
- ✅ Backward compatibility maintained (old constructors still work)
- ✅ Storage errors don't prevent MQTT operations
- ✅ All tests pass with >90% coverage
- ✅ Mock store allows easy testing

## Migration Notes
- Existing code using old constructors continues to work
- Storage parameter is optional (nil = no persistence)
- Log warnings for storage errors, don't fail operations
- This allows gradual migration and fallback behavior
