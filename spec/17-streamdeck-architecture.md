# Stream Deck+ Interface Architecture

## Overview

This specification defines the architecture for a Stream Deck+ interface that provides tactile control of office lights through buttons, rotary encoders, and a touchscreen display.

## Hardware Components

### Stream Deck+ Features

1. **Buttons**: 8 LCD buttons (2 rows × 4 columns)
   - Each button: 120×120 pixel display
   - Programmable images
   - Touch-sensitive

2. **Touchscreen**: Long horizontal strip (800×100 pixels)
   - Touch input only (no multi-touch or gestures)
   - Single tap detection

3. **Rotary Encoders**: 4 physical dials
   - Detented rotation (discrete up/down events)
   - Click/press capability
   - Positioned below touchscreen

## Interface Design

### Mode Selection (Second Row Buttons)

The 4 buttons in the second row act as a radio button set:

| Button 5 | Button 6 | Button 7 | Button 8 |
|----------|----------|----------|----------|
| LED Strip | LED Bar RGBW | LED Bar White | Video Lights |

**Visual States:**
- **Active**: Bright/highlighted icon
- **Inactive**: Dimmed icon

**Icons:**
- LED Strip: RGB gradient bar
- LED Bar RGBW: Multi-color segment display
- LED Bar White: White light bulb
- Video Lights: Camera/studio light icon

### Touchscreen Display (800×100)

Divided into 4 equal sections (200×100 each), showing values based on selected mode:

#### Mode 1: LED Strip (3 active sections)
```
┌────────┬────────┬────────┬────────┐
│  Red   │ Green  │  Blue  │ (empty)│
│  ###   │  ###   │  ###   │        │
└────────┴────────┴────────┴────────┘
```

#### Mode 2: LED Bar RGBW (4 active sections)
```
┌────────┬────────┬────────┬────────┐
│  Red   │ Green  │  Blue  │ White  │
│  ###   │  ###   │  ###   │  ###   │
└────────┴────────┴────────┴────────┘
```

#### Mode 3: LED Bar White (2 active sections)
```
┌────────┬────────┬────────┬────────┐
│Section1│Section2│ (empty)│ (empty)│
│  ###   │  ###   │        │        │
└────────┴────────┴────────┴────────┘
```

#### Mode 4: Video Lights (2 active sections)
```
┌────────┬────────┬────────┬────────┐
│ Light1 │ Light2 │ (empty)│ (empty)│
│  ###   │  ###   │        │        │
│ ON/OFF │ ON/OFF │        │        │
└────────┴────────┴────────┴────────┘
```

**Display Elements per Section:**
- Label (top): Component name
- Value (center): Numeric value (0-255 or 0-100)
- Bar graph (optional): Visual indicator

### Rotary Encoders (4 Dials)

**Behavior:**
- **Rotation**: Increment/decrement the corresponding touchscreen value
  - Clockwise: Increase by 5 (or 1 with modifier?)
  - Counter-clockwise: Decrease by 5

- **Click/Press**: Toggle behavior
  - LED Strip/LED Bar: Toggle between 0 and last non-zero value
  - Video Lights: Toggle on/off state

**Dial Mapping:**
- Dial 1 → Touchscreen Section 1
- Dial 2 → Touchscreen Section 2
- Dial 3 → Touchscreen Section 3
- Dial 4 → Touchscreen Section 4

**Inactive Dials:**
- When fewer than 4 values are shown (e.g., LED Strip has only 3), the unused dials are inactive

## Technical Architecture

### Stream Deck SDK Integration

**Available Options:**
1. **Elgato Stream Deck SDK** (Official)
   - WebSocket-based communication
   - Plugin architecture
   - Requires manifest.json and JavaScript/HTML

2. **Go Libraries:**
   - `github.com/muesli/streamdeck` - Direct USB HID communication
   - `github.com/magicmonkey/go-streamdeck` - Alternative Go library

**Recommended Approach:** Use `github.com/muesli/streamdeck`
- Direct control without Stream Deck software dependency
- Pure Go implementation
- Supports all Stream Deck+ features
- Better integration with our existing Go codebase

### Component Architecture

```
streamdeck/
├── streamdeck.go       # Main interface controller
├── model.go            # State model and mode enum
├── render.go           # Image rendering for buttons/touchscreen
├── events.go           # Event handling (buttons, dials, touch)
├── modes.go            # Mode-specific logic
├── icons/              # Button icons (120x120 PNG)
│   ├── ledstrip.png
│   ├── ledbar_rgbw.png
│   ├── ledbar_white.png
│   └── videolight.png
└── fonts/              # Fonts for touchscreen rendering
    └── Roboto-Regular.ttf
```

### State Management

```go
type Mode int

const (
    ModeLEDStrip Mode = iota
    ModeLEDBarRGBW
    ModeLEDBarWhite
    ModeVideoLights
)

type StreamDeckUI struct {
    device      *streamdeck.Device
    ledStrip    *ledstrip.LEDStrip
    ledBar      *ledbar.LEDBar
    videoLight1 *videolight.VideoLight
    videoLight2 *videolight.VideoLight

    currentMode Mode
    lastValues  [4]int // Store last non-zero values for toggle

    // Rendering
    buttonImages [8]image.Image
    touchImage   image.Image
}
```

### Event Handling

**Button Events:**
```go
func (s *StreamDeckUI) handleButtonPress(buttonIndex int)
```
- Buttons 0-3 (top row): Reserved for future functionality
- Buttons 4-7 (second row): Mode selection

**Dial Events:**
```go
func (s *StreamDeckUI) handleDialRotate(dialIndex int, ticks int)
func (s *StreamDeckUI) handleDialPress(dialIndex int)
```

**Touch Events:**
```go
func (s *StreamDeckUI) handleTouch(x, y int)
```
- Calculate which section was touched (x / 200)
- Provide visual feedback (optional)

### Rendering Pipeline

**Button Rendering:**
1. Load base icon from `icons/` directory
2. If active mode, apply highlight effect
3. Encode as JPEG
4. Send to Stream Deck button

**Touchscreen Rendering:**
1. Create 800×100 image
2. Divide into 4 sections (200×100 each)
3. For each active section:
   - Draw background
   - Draw label text
   - Draw value text
   - Draw bar graph (optional)
4. Encode as JPEG
5. Send to Stream Deck touchscreen

### Integration with Drivers

**LED Strip Mode:**
```go
// Read current values
r, g, b := s.ledStrip.GetColor()

// Update value via dial
s.ledStrip.SetColor(r, g, b) // Publishes to MQTT
```

**LED Bar RGBW Mode:**
- Show values for LED index 0 in selected section
- Use `SetRGBW()` to update individual LED
- **Decision needed:** Which LED index to control? First LED (0)? Average of all 6?

**LED Bar White Mode:**
- Show average brightness for all 13 white LEDs in each section
- Update all 13 LEDs together using `SetAllWhite()` per section

**Video Lights Mode:**
```go
// Read states
on1, brightness1 := s.videoLight1.GetState()
on2, brightness2 := s.videoLight2.GetState()

// Update via dial
s.videoLight1.TurnOn(brightness)
s.videoLight1.TurnOff()
```

## Polling vs. Event-Driven Updates

**Touchscreen Updates:**
- Poll driver state every 100ms
- Redraw only if values changed
- Reduces unnecessary rendering

**Button Updates:**
- Update immediately on mode change
- No polling needed

## Error Handling

**Device Connection:**
- Detect if Stream Deck is connected
- Graceful failure if device not found
- Auto-reconnect on disconnect (optional)

**Rendering Errors:**
- Log errors but don't crash
- Show error state on touchscreen

## Concurrency

- Run in separate goroutine (like TUI and Web)
- Use mutex-protected driver access (drivers already thread-safe via MQTT client)
- Event loop for Stream Deck input
- Separate goroutine for periodic touchscreen updates

## Performance Considerations

**Image Rendering:**
- Pre-render mode button images (cache)
- Only re-render touchscreen when values change
- Use efficient image encoding (JPEG quality 90%)

**Update Rate:**
- Touchscreen: 100ms (10 FPS)
- Buttons: On-demand only
- Dials: Immediate response

## Future Enhancements

1. **Top Row Buttons (0-3):**
   - Preset scenes
   - Quick toggles
   - Brightness presets

2. **Touchscreen Gestures:**
   - If hardware supports: Swipe to change modes

3. **Visual Feedback:**
   - Animations on value change
   - Color preview for LED strip/bar

4. **Configuration:**
   - Customizable dial sensitivity
   - Custom button icons

## Dependencies

```bash
go get github.com/muesli/streamdeck
go get golang.org/x/image/font
go get golang.org/x/image/font/basicfont
go get golang.org/x/image/draw
```

## Platform Considerations

**Linux:**
- Requires udev rules for USB access
- Add user to `plugdev` group

**macOS:**
- May require sudo or entitlements
- Test with standard user permissions

**Windows:**
- Should work out of the box
- May need USB driver

## Testing Strategy

**Without Hardware:**
- Mock Stream Deck interface
- Test rendering to PNG files
- Test state management logic

**With Hardware:**
- Manual testing of all modes
- Verify dial sensitivity
- Test concurrent UI operation

---

## Summary

The Stream Deck+ interface provides a tactile, visual control surface for office lights with:
- 4 operational modes via button selection
- Real-time value display on touchscreen
- Physical dial control with toggle functionality
- Concurrent operation with TUI and Web interfaces
- Direct USB HID communication
- Efficient rendering and polling

Next: See `18-streamdeck-implementation-plan.md` for step-by-step implementation guide.
