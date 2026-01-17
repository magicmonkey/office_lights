# Usage Examples

This guide shows how to use the office lights control system drivers programmatically.

## Basic Setup

```go
package main

import (
    "log"

    "github.com/kevin/office_lights/drivers/ledstrip"
    "github.com/kevin/office_lights/drivers/ledbar"
    "github.com/kevin/office_lights/drivers/videolight"
    "github.com/kevin/office_lights/mqtt"
)

func main() {
    // Create MQTT client
    config := mqtt.Config{
        Broker:   "tcp://localhost:1883",
        ClientID: "my_lights_controller",
    }

    client, err := mqtt.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }

    if err := client.Connect(); err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect()

    // Create driver instances
    strip := ledstrip.NewLEDStrip(client, mqtt.TopicLEDStrip)
    bar, _ := ledbar.NewLEDBar(0, client, mqtt.TopicLEDBar)
    light1, _ := videolight.NewVideoLight(1, client, mqtt.TopicVideoLight1)
    light2, _ := videolight.NewVideoLight(2, client, mqtt.TopicVideoLight2)

    // Use the lights (examples below)
}
```

## LED Strip Examples

### Set a specific color
```go
// Set to warm white
strip.SetColor(255, 200, 150)

// Set to bright red
strip.SetColor(255, 0, 0)

// Set to cyan
strip.SetColor(0, 255, 255)
```

### Use preset colors
```go
strip.SetRed()      // Pure red
strip.SetGreen()    // Pure green
strip.SetBlue()     // Pure blue
strip.SetWhite()    // White
strip.SetYellow()   // Yellow
strip.SetCyan()     // Cyan
strip.SetMagenta()  // Magenta
```

### Adjust brightness
```go
// Set a color first
strip.SetColor(200, 150, 100)

// Then adjust brightness to 50%
strip.SetBrightness(50)  // Result: (100, 75, 50)
```

### Turn off
```go
strip.TurnOff()  // Sets to (0, 0, 0)
```

### Check current state
```go
r, g, b := strip.GetColor()
log.Printf("Current color: R=%d, G=%d, B=%d", r, g, b)
```

## Video Light Examples

### Basic on/off control
```go
// Turn on at 75% brightness
light1.TurnOn(75)

// Turn off
light1.TurnOff()
```

### Set state directly
```go
// Turn on at specific brightness
light1.SetState(true, 50)

// Turn off but remember brightness
light1.SetState(false, 50)
```

### Adjust brightness while on
```go
// First turn on
light1.TurnOn(50)

// Later adjust brightness
light1.SetBrightness(75)  // Still on, but now at 75%
```

### Check current state
```go
on, brightness := light1.GetState()
log.Printf("Light 1: on=%v, brightness=%d", on, brightness)
```

### Control multiple lights together
```go
// Turn on both lights at same brightness
light1.TurnOn(80)
light2.TurnOn(80)

// Or turn them off together
light1.TurnOff()
light2.TurnOff()
```

## LED Bar Examples

### Set individual RGBW LEDs
```go
// Set first RGBW LED in section 1 to blue
bar.SetRGBW(1, 0, 0, 0, 255, 100)  // section, index, r, g, b, w

// Set last RGBW LED in section 2 to white
bar.SetRGBW(2, 5, 255, 255, 255, 255)
```

### Set individual white LEDs
```go
// Set first white LED in section 1
bar.SetWhite(1, 0, 200)  // section, index, value

// Set all white LEDs in section 2 individually
for i := 0; i < 13; i++ {
    bar.SetWhite(2, i, 150)
}
```

### Set all LEDs at once
```go
// Set all RGBW LEDs to warm white
bar.SetAllRGBW(255, 200, 150, 100)

// Set all white LEDs to same brightness
bar.SetAllWhite(180)
```

### Turn off sections
```go
// Turn off just section 1
bar.TurnOffSection(1)

// Turn off just section 2
bar.TurnOffSection(2)

// Turn off entire bar
bar.TurnOffAll()
```

### Check current state
```go
// Get RGBW LED value
r, g, b, w, err := bar.GetRGBW(1, 0)
if err == nil {
    log.Printf("RGBW LED: R=%d, G=%d, B=%d, W=%d", r, g, b, w)
}

// Get white LED value
value, err := bar.GetWhite(1, 0)
if err == nil {
    log.Printf("White LED: %d", value)
}
```

## Complete Example: Scene Control

```go
// Create a "Video Call" scene
func VideoCallScene(strip *ledstrip.LEDStrip, bar *ledbar.LEDBar,
                    light1, light2 *videolight.VideoLight) {
    // Bright white video lights
    light1.TurnOn(100)
    light2.TurnOn(100)

    // Warm ambient strip
    strip.SetColor(255, 200, 150)
    strip.SetBrightness(30)

    // Turn off LED bar to avoid color conflicts
    bar.TurnOffAll()
}

// Create a "Focus" scene
func FocusScene(strip *ledstrip.LEDStrip, bar *ledbar.LEDBar,
                light1, light2 *videolight.VideoLight) {
    // Moderate video lights
    light1.TurnOn(60)
    light2.TurnOn(60)

    // Cool white bar for desk lighting
    bar.SetAllRGBW(255, 255, 255, 200)
    bar.SetAllWhite(180)

    // Subtle blue ambient
    strip.SetColor(100, 150, 255)
    strip.SetBrightness(20)
}

// Create an "Off" scene
func AllOffScene(strip *ledstrip.LEDStrip, bar *ledbar.LEDBar,
                 light1, light2 *videolight.VideoLight) {
    strip.TurnOff()
    bar.TurnOffAll()
    light1.TurnOff()
    light2.TurnOff()
}

// Create a "Party" scene
func PartyScene(strip *ledstrip.LEDStrip, bar *ledbar.LEDBar,
                light1, light2 *videolight.VideoLight) {
    // Colorful strip
    strip.SetMagenta()

    // Rainbow on bar - alternate colors
    bar.SetRGBW(1, 0, 255, 0, 0, 0)      // Red
    bar.SetRGBW(1, 1, 255, 127, 0, 0)    // Orange
    bar.SetRGBW(1, 2, 255, 255, 0, 0)    // Yellow
    bar.SetRGBW(1, 3, 0, 255, 0, 0)      // Green
    bar.SetRGBW(1, 4, 0, 0, 255, 0)      // Blue
    bar.SetRGBW(1, 5, 148, 0, 211, 0)    // Purple

    // Dim video lights for ambiance
    light1.TurnOn(20)
    light2.TurnOn(20)
}

// Usage in main
func main() {
    // ... setup code from above ...

    // Apply a scene
    VideoCallScene(strip, bar, light1, light2)

    // Later switch to another scene
    FocusScene(strip, bar, light1, light2)
}
```

## Error Handling

All driver methods that can fail return an error. Always check errors:

```go
// Good - check errors
if err := strip.SetColor(100, 150, 200); err != nil {
    log.Printf("Failed to set strip color: %v", err)
}

// Also good - handle validation errors
if err := bar.SetRGBW(1, 0, 300, 0, 0, 0); err != nil {
    log.Printf("Invalid RGB value: %v", err)  // Will fail - 300 > 255
}

// Good - check constructor errors
light, err := videolight.NewVideoLight(0, client, topic)
if err != nil {
    log.Fatal("Invalid light ID")  // Will fail - ID must be >= 1
}
```

## Input Validation

All drivers validate inputs:

### LED Strip
- RGB values must be 0-255
- Brightness must be 0-100

### LED Bar
- Section must be 1 or 2
- RGBW LED index must be 0-5
- White LED index must be 0-12
- All color/brightness values must be 0-255

### Video Light
- Light ID must be >= 1
- Brightness must be 0-100

Invalid inputs will return descriptive errors without publishing any MQTT messages.

## Testing

When testing your code, use the mock publisher:

```go
import "github.com/kevin/office_lights/mqtt"

func TestMyLightControl(t *testing.T) {
    // Create mock instead of real MQTT client
    mock := mqtt.NewMockPublisher()

    // Create driver with mock
    strip := ledstrip.NewLEDStrip(mock, "test/topic")

    // Use the driver
    strip.SetColor(100, 150, 200)

    // Verify message was published
    msg := mock.GetLastMessage()
    if msg == nil {
        t.Fatal("No message published")
    }

    // Check message count
    if mock.MessageCount() != 1 {
        t.Errorf("Expected 1 message, got %d", mock.MessageCount())
    }

    // Clear for next test
    mock.Clear()
}
```
