# LED Bar Driver

## Overview
Implement the driver for RGBW LED bars with white-only LEDs via MQTT.

## Tasks

### 1. Create Package Structure
- Create `drivers/ledbar/ledbar.go`
- Define package `ledbar`

### 2. Define State Structure
```go
type LEDBar struct {
    RGBW1 [6][4]int  // First set of 6 RGBW LEDs
    White1 [13]int   // First set of 13 white LEDs
    RGBW2 [6][4]int  // Second set of 6 RGBW LEDs
    White2 [13]int   // Second set of 13 white LEDs
    BarID int        // Bar identifier (0 or 1)
    // MQTT publisher reference
}
```

### 3. Understand Message Format
The comma-separated message has this structure:
- Values 0-23: 6 RGBW LEDs (4 values each: R,G,B,W)
- Values 24-36: 13 white LEDs (1 value each)
- Values 37-39: 3 ignored values
- Values 40-63: 6 RGBW LEDs (4 values each: R,G,B,W)
- Values 64-76: 13 white LEDs (1 value each)

Total: 77 comma-separated values

### 4. Implement Constructor
- `NewLEDBar(barID int, mqttClient interface{}) *LEDBar`
- Initialize with default state (all LEDs off)
- Validate barID

### 5. Implement State Management Methods
- `SetRGBW(section int, index int, r, g, b, w int) error`
  - section: 1 or 2 (first or second group)
  - index: 0-5 (which RGBW LED)
  - Validate all inputs
- `SetWhite(section int, index int, value int) error`
  - section: 1 or 2
  - index: 0-12
  - Validate inputs

### 6. Implement MQTT Message Formatting
- `formatMessage() string`
  - Build comma-separated string with all 77 values
  - Order: RGBW1, White1, ignored(3x), RGBW2, White2
  - Ensure correct value count

### 7. Implement Publish Method
- `Publish() error`
  - Format message
  - Publish to topic `kevinoffice/ledbar/<barID>`
  - Handle errors

### 8. Add Convenience Methods
- `TurnOffSection(section int)` - Turn off one section
- `TurnOffAll()` - Turn off entire bar
- `SetAllRGBW(r, g, b, w int)` - Set all RGBW LEDs to same color
- `SetAllWhite(value int)` - Set all white LEDs to same value

## Success Criteria
- Can create LED bar instance for each bar ID
- Can set individual RGBW and white LED values
- Generates correct comma-separated message (77 values)
- Publishes to correct MQTT topic
- Validates all input ranges
