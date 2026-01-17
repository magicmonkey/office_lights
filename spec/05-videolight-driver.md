# Video Light Driver

## Overview
Implement the driver for video lights via MQTT.

## Tasks

### 1. Create Package Structure
- Create `drivers/videolight/videolight.go`
- Define package `videolight`

### 2. Define State Structure
```go
type VideoLight struct {
    On bool        // Light on/off state
    Brightness int // Brightness level (0-100)
    LightID int    // Light identifier (1 or 2)
    // MQTT publisher reference
}
```

### 3. Implement Constructor
- `NewVideoLight(lightID int, mqttClient interface{}) *VideoLight`
- Initialize with default state (off, brightness 0)
- Validate lightID (should be 1 or 2)

### 4. Implement State Management Methods
- `SetState(on bool, brightness int) error`
  - Validate brightness (0-100 range)
  - Update internal state
  - Trigger MQTT message send
- `TurnOn(brightness int) error`
  - Convenience method
- `TurnOff() error`
  - Convenience method
- `SetBrightness(brightness int) error`
  - Set brightness while maintaining on/off state
- `GetState() (bool, int)`
  - Return current state

### 5. Implement MQTT Message Formatting
- `formatMessage() string`
  - Create message: `set,<on>,<brightness>`
  - Example: `set,true,50`
  - Convert boolean to lowercase string ("true"/"false")
  - Return formatted string

### 6. Implement Publish Method
- `Publish() error`
  - Format message
  - Publish to topic `kevinoffice/videolight/<lightID>/command/light:0`
  - Handle errors

### 7. Input Validation
- Ensure brightness is 0-100
- Ensure lightID is valid
- Handle edge cases (e.g., turning off should set appropriate brightness)

## Success Criteria
- Can create video light instance for each light ID
- Can set on/off state and brightness
- Generates correct message format: `set,<bool>,<int>`
- Publishes to correct MQTT topic for each light
- Validates input ranges
