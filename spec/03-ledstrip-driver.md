# LED Strip Driver

## Overview
Implement the driver for RGB LED strip control via MQTT.

## Tasks

### 1. Create Package Structure
- Create `drivers/ledstrip/ledstrip.go`
- Define package `ledstrip`

### 2. Define State Structure
```go
type LEDStrip struct {
    R int // Red value (0-255)
    G int // Green value (0-255)
    B int // Blue value (0-255)
    // MQTT publisher reference
}
```

### 3. Implement Constructor
- `NewLEDStrip(mqttClient interface{}) *LEDStrip`
- Initialize with default state (e.g., off = 0,0,0)

### 4. Implement State Management Methods
- `SetColor(r, g, b int) error`
  - Validate RGB values (0-255 range)
  - Update internal state
  - Trigger MQTT message send
- `GetColor() (int, int, int)`
  - Return current RGB values

### 5. Implement MQTT Message Formatting
- `formatMessage() ([]byte, error)`
  - Create JSON structure: `{"sequence":"fill", "data":{"r":<int>,"g":<int>,"b":<int>}}`
  - Use `encoding/json` to marshal
  - Return byte array for MQTT publish

### 6. Implement Publish Method
- `Publish() error`
  - Format message
  - Publish to topic `kevinoffice/ledstrip/sequence`
  - Handle errors

### 7. Add Convenience Methods (Optional)
- `TurnOff()` - Set to 0,0,0
- `SetBrightness(percentage int)` - Scale current color
- Preset colors (red, green, blue, white, etc.)

## Success Criteria
- Can create LED strip instance
- Can set RGB values
- Generates correct JSON message format
- Publishes to correct MQTT topic
- Validates input ranges
