# Stream Deck+ Interface Implementation Plan

## Overview

This document provides a step-by-step guide for implementing the Stream Deck+ interface for office lights control.

**Estimated Complexity:** High
**Prerequisites:** Stream Deck+ hardware (optional for initial development)

## Implementation Phases

### Phase 1: Dependencies and Project Setup

**Goal:** Install required libraries and create package structure

**Steps:**

1. Install Stream Deck library:
```bash
go get github.com/muesli/streamdeck
```

2. Install image processing libraries:
```bash
go get golang.org/x/image/font
go get golang.org/x/image/font/basicfont
go get golang.org/x/image/math/fixed
go get golang.org/x/image/draw
```

3. Create package structure:
```bash
mkdir -p streamdeck/icons
mkdir -p streamdeck/fonts
```

4. Create initial files:
   - `streamdeck/streamdeck.go` - Main entry point
   - `streamdeck/model.go` - State structures
   - `streamdeck/render.go` - Image rendering
   - `streamdeck/events.go` - Event handling
   - `streamdeck/modes.go` - Mode-specific logic

**Validation:**
- All dependencies installed without errors
- Package structure created
- Files compile (even if empty)

---

### Phase 2: Core Data Structures

**Goal:** Define the state model and mode enumeration

**File:** `streamdeck/model.go`

```go
package streamdeck

import (
	"image"
	"sync"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
	"github.com/muesli/streamdeck"
)

// Mode represents the current control mode
type Mode int

const (
	ModeLEDStrip Mode = iota
	ModeLEDBarRGBW
	ModeLEDBarWhite
	ModeVideoLights
)

// String returns the mode name
func (m Mode) String() string {
	switch m {
	case ModeLEDStrip:
		return "LED Strip"
	case ModeLEDBarRGBW:
		return "LED Bar RGBW"
	case ModeLEDBarWhite:
		return "LED Bar White"
	case ModeVideoLights:
		return "Video Lights"
	default:
		return "Unknown"
	}
}

// StreamDeckUI manages the Stream Deck+ interface
type StreamDeckUI struct {
	device      *streamdeck.Device
	ledStrip    *ledstrip.LEDStrip
	ledBar      *ledbar.LEDBar
	videoLight1 *videolight.VideoLight
	videoLight2 *videolight.VideoLight

	mu          sync.Mutex
	currentMode Mode
	lastValues  [4]int // Store last non-zero values for dial toggle

	// For graceful shutdown
	quit chan struct{}
}

// Section represents a touchscreen section
type Section struct {
	Label string
	Value int
	Max   int
}
```

**Validation:**
- Code compiles
- Enums defined correctly
- Struct fields match requirements

---

### Phase 3: Device Initialization

**Goal:** Detect and initialize Stream Deck+ device

**File:** `streamdeck/streamdeck.go`

```go
package streamdeck

import (
	"fmt"
	"log"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
	"github.com/muesli/streamdeck"
)

// New creates a new Stream Deck UI instance
func New(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) (*StreamDeckUI, error) {
	// Find Stream Deck device
	devices, err := streamdeck.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate devices: %w", err)
	}

	if len(devices) == 0 {
		return nil, fmt.Errorf("no Stream Deck device found")
	}

	// Open first device (assume Stream Deck+)
	dev, err := streamdeck.Open(devices[0])
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %w", err)
	}

	// Reset device to clean state
	if err := dev.Reset(); err != nil {
		dev.Close()
		return nil, fmt.Errorf("failed to reset device: %w", err)
	}

	// Set brightness
	if err := dev.SetBrightness(80); err != nil {
		log.Printf("Warning: failed to set brightness: %v", err)
	}

	ui := &StreamDeckUI{
		device:      dev,
		ledStrip:    strip,
		ledBar:      bar,
		videoLight1: vl1,
		videoLight2: vl2,
		currentMode: ModeLEDStrip,
		quit:        make(chan struct{}),
	}

	return ui, nil
}

// Close cleans up resources
func (s *StreamDeckUI) Close() error {
	close(s.quit)
	if s.device != nil {
		s.device.Reset()
		return s.device.Close()
	}
	return nil
}

// Run starts the Stream Deck UI event loop
func Run(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) error {
	ui, err := New(strip, bar, vl1, vl2)
	if err != nil {
		return err
	}
	defer ui.Close()

	// Initial render
	if err := ui.renderAll(); err != nil {
		return fmt.Errorf("initial render failed: %w", err)
	}

	// Start event loop
	return ui.eventLoop()
}
```

**Validation:**
- Device detection works (if hardware available)
- Graceful error if device not found
- Device resets to clean state
- Brightness set successfully

---

### Phase 4: Image Rendering Foundation

**Goal:** Implement button and touchscreen rendering

**File:** `streamdeck/render.go`

```go
package streamdeck

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// renderAll renders all buttons and touchscreen
func (s *StreamDeckUI) renderAll() error {
	if err := s.renderButtons(); err != nil {
		return err
	}
	if err := s.renderTouchscreen(); err != nil {
		return err
	}
	return nil
}

// renderButtons renders all mode selection buttons
func (s *StreamDeckUI) renderButtons() error {
	// Clear top row buttons (reserved for future functionality)
	for i := 0; i < 4; i++ {
		img := s.createEmptyButton()
		s.setButtonImage(i, img)
	}

	// Render mode selection buttons on second row (4-7)
	modes := []Mode{ModeLEDStrip, ModeLEDBarRGBW, ModeLEDBarWhite, ModeVideoLights}
	for i, mode := range modes {
		img := s.createButtonImage(mode, mode == s.currentMode)
		if err := s.setButtonImage(i+4, img); err != nil {
			log.Printf("Failed to set button %d: %v", i+4, err)
		}
	}

	return nil
}

// createButtonImage creates a 120x120 button image
func (s *StreamDeckUI) createButtonImage(mode Mode, active bool) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 120, 120))

	// Background color
	bgColor := color.RGBA{40, 40, 40, 255}
	if active {
		bgColor = color.RGBA{80, 120, 200, 255}
	}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// Draw text label
	label := mode.String()
	s.drawCenteredText(img, label, 60, color.White)

	return img
}

// createEmptyButton creates a blank button
func (s *StreamDeckUI) createEmptyButton() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 120, 120))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	return img
}

// renderTouchscreen renders the touchscreen display
func (s *StreamDeckUI) renderTouchscreen() error {
	img := s.createTouchscreenImage()
	return s.setTouchscreenImage(img)
}

// createTouchscreenImage creates 800x100 touchscreen image
func (s *StreamDeckUI) createTouchscreenImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 800, 100))

	// Black background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)

	// Get sections for current mode
	sections := s.getSections()

	// Draw each section (200px wide each)
	for i, section := range sections {
		x := i * 200
		if section.Label != "" {
			s.drawSection(img, x, section)
		}
	}

	return img
}

// drawSection draws a single touchscreen section
func (s *StreamDeckUI) drawSection(img *image.RGBA, x int, section Section) {
	// Section border
	for i := 0; i < 200; i++ {
		for j := 0; j < 100; j++ {
			px := x + i
			if i == 0 || i == 199 {
				img.Set(px, j, color.RGBA{100, 100, 100, 255})
			}
		}
	}

	// Draw label (top)
	s.drawText(img, section.Label, x+10, 20, color.RGBA{150, 150, 150, 255})

	// Draw value (center)
	valueStr := fmt.Sprintf("%d", section.Value)
	s.drawText(img, valueStr, x+10, 60, color.White)
}

// setButtonImage sends image to a button
func (s *StreamDeckUI) setButtonImage(index int, img image.Image) error {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return err
	}
	return s.device.SetImage(index, buf.Bytes())
}

// setTouchscreenImage sends image to touchscreen
func (s *StreamDeckUI) setTouchscreenImage(img image.Image) error {
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		return err
	}
	return s.device.SetLCD(buf.Bytes())
}

// drawText draws text at position
func (s *StreamDeckUI) drawText(img *image.RGBA, text string, x, y int, col color.Color) {
	point := fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)
}

// drawCenteredText draws centered text
func (s *StreamDeckUI) drawCenteredText(img *image.RGBA, text string, y int, col color.Color) {
	// Simple centering (approximate)
	textWidth := len(text) * 7 // basicfont is ~7px wide
	x := (120 - textWidth) / 2
	s.drawText(img, text, x, y, col)
}
```

**Validation:**
- Buttons render with correct labels
- Active button has different color
- Touchscreen renders 4 sections
- Text is readable

---

### Phase 5: Mode-Specific Section Data

**Goal:** Implement logic to get section data based on current mode

**File:** `streamdeck/modes.go`

```go
package streamdeck

import "fmt"

// getSections returns touchscreen sections for current mode
func (s *StreamDeckUI) getSections() []Section {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch s.currentMode {
	case ModeLEDStrip:
		return s.getLEDStripSections()
	case ModeLEDBarRGBW:
		return s.getLEDBarRGBWSections()
	case ModeLEDBarWhite:
		return s.getLEDBarWhiteSections()
	case ModeVideoLights:
		return s.getVideoLightsSections()
	default:
		return make([]Section, 4)
	}
}

// getLEDStripSections returns RGB sections
func (s *StreamDeckUI) getLEDStripSections() []Section {
	r, g, b := s.ledStrip.GetColor()
	return []Section{
		{Label: "Red", Value: r, Max: 255},
		{Label: "Green", Value: g, Max: 255},
		{Label: "Blue", Value: b, Max: 255},
		{Label: "", Value: 0, Max: 0}, // Empty
	}
}

// getLEDBarRGBWSections returns RGBW sections for LED 0
func (s *StreamDeckUI) getLEDBarRGBWSections() []Section {
	// Use first LED (index 0) of first section as representative
	r, g, b, w, err := s.ledBar.GetRGBW(1, 0)
	if err != nil {
		log.Printf("Error getting LED bar RGBW: %v", err)
		return make([]Section, 4)
	}

	return []Section{
		{Label: "Red", Value: r, Max: 255},
		{Label: "Green", Value: g, Max: 255},
		{Label: "Blue", Value: b, Max: 255},
		{Label: "White", Value: w, Max: 255},
	}
}

// getLEDBarWhiteSections returns average white brightness per section
func (s *StreamDeckUI) getLEDBarWhiteSections() []Section {
	// Calculate average for section 1
	sum1 := 0
	for i := 0; i < 13; i++ {
		val, err := s.ledBar.GetWhite(1, i)
		if err == nil {
			sum1 += val
		}
	}
	avg1 := sum1 / 13

	// Calculate average for section 2
	sum2 := 0
	for i := 0; i < 13; i++ {
		val, err := s.ledBar.GetWhite(2, i)
		if err == nil {
			sum2 += val
		}
	}
	avg2 := sum2 / 13

	return []Section{
		{Label: "Section 1", Value: avg1, Max: 255},
		{Label: "Section 2", Value: avg2, Max: 255},
		{Label: "", Value: 0, Max: 0}, // Empty
		{Label: "", Value: 0, Max: 0}, // Empty
	}
}

// getVideoLightsSections returns video light brightness
func (s *StreamDeckUI) getVideoLightsSections() []Section {
	on1, brightness1 := s.videoLight1.GetState()
	on2, brightness2 := s.videoLight2.GetState()

	label1 := "Light 1: OFF"
	label2 := "Light 2: OFF"
	if on1 {
		label1 = "Light 1: ON"
	}
	if on2 {
		label2 = "Light 2: ON"
	}

	return []Section{
		{Label: label1, Value: brightness1, Max: 100},
		{Label: label2, Value: brightness2, Max: 100},
		{Label: "", Value: 0, Max: 0}, // Empty
		{Label: "", Value: 0, Max: 0}, // Empty
	}
}
```

**Validation:**
- Each mode returns correct section data
- Empty sections have empty labels
- Values match driver state

---

### Phase 6: Event Handling

**Goal:** Handle button presses, dial rotations, and dial clicks

**File:** `streamdeck/events.go`

```go
package streamdeck

import (
	"log"
	"time"
)

// eventLoop processes Stream Deck events
func (s *StreamDeckUI) eventLoop() error {
	// Create ticker for periodic touchscreen updates
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Button press channel
	btnChan := make(chan int, 10)
	s.device.SetButtonPressCallback(func(btnIndex int, pressed bool) {
		if pressed {
			btnChan <- btnIndex
		}
	})

	// Dial rotation channel
	dialRotateChan := make(chan struct{ index, ticks int }, 10)
	s.device.SetDialRotateCallback(func(dialIndex, ticks int) {
		dialRotateChan <- struct{ index, ticks int }{dialIndex, ticks}
	})

	// Dial press channel
	dialPressChan := make(chan int, 10)
	s.device.SetDialPressCallback(func(dialIndex int, pressed bool) {
		if pressed {
			dialPressChan <- dialIndex
		}
	})

	for {
		select {
		case <-s.quit:
			return nil

		case btnIndex := <-btnChan:
			s.handleButtonPress(btnIndex)

		case dr := <-dialRotateChan:
			s.handleDialRotate(dr.index, dr.ticks)

		case dialIndex := <-dialPressChan:
			s.handleDialPress(dialIndex)

		case <-ticker.C:
			// Periodic touchscreen update
			if err := s.renderTouchscreen(); err != nil {
				log.Printf("Failed to update touchscreen: %v", err)
			}
		}
	}
}

// handleButtonPress processes button press events
func (s *StreamDeckUI) handleButtonPress(btnIndex int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Top row buttons (0-3): Reserved for future functionality
	if btnIndex < 4 {
		return
	}

	// Second row buttons (4-7) select mode
	newMode := Mode(btnIndex - 4)
	if newMode != s.currentMode {
		s.currentMode = newMode
		log.Printf("Switched to mode: %s", newMode)

		// Re-render buttons to show active state
		go s.renderButtons()
	}
}

// handleDialRotate processes dial rotation
func (s *StreamDeckUI) handleDialRotate(dialIndex, ticks int) {
	s.mu.Lock()
	mode := s.currentMode
	s.mu.Unlock()

	// Calculate value change (5 per tick)
	delta := ticks * 5

	switch mode {
	case ModeLEDStrip:
		s.rotateLEDStrip(dialIndex, delta)
	case ModeLEDBarRGBW:
		s.rotateLEDBarRGBW(dialIndex, delta)
	case ModeLEDBarWhite:
		s.rotateLEDBarWhite(dialIndex, delta)
	case ModeVideoLights:
		s.rotateVideoLights(dialIndex, delta)
	}
}

// handleDialPress processes dial click
func (s *StreamDeckUI) handleDialPress(dialIndex int) {
	s.mu.Lock()
	mode := s.currentMode
	s.mu.Unlock()

	switch mode {
	case ModeLEDStrip:
		s.toggleLEDStrip(dialIndex)
	case ModeLEDBarRGBW:
		s.toggleLEDBarRGBW(dialIndex)
	case ModeLEDBarWhite:
		s.toggleLEDBarWhite(dialIndex)
	case ModeVideoLights:
		s.toggleVideoLights(dialIndex)
	}
}

// Rotation handlers for each mode
func (s *StreamDeckUI) rotateLEDStrip(dialIndex, delta int) {
	r, g, b := s.ledStrip.GetColor()

	switch dialIndex {
	case 0: // Red
		r = clamp(r+delta, 0, 255)
	case 1: // Green
		g = clamp(g+delta, 0, 255)
	case 2: // Blue
		b = clamp(b+delta, 0, 255)
	default:
		return // Dial 3 unused
	}

	s.ledStrip.SetColor(r, g, b)
}

func (s *StreamDeckUI) rotateLEDBarRGBW(dialIndex, delta int) {
	if dialIndex > 3 {
		return
	}

	r, g, b, w, _ := s.ledBar.GetRGBW(1, 0)

	switch dialIndex {
	case 0:
		r = clamp(r+delta, 0, 255)
	case 1:
		g = clamp(g+delta, 0, 255)
	case 2:
		b = clamp(b+delta, 0, 255)
	case 3:
		w = clamp(w+delta, 0, 255)
	}

	// Update all RGBW LEDs in both sections
	for i := 0; i < 6; i++ {
		s.ledBar.SetRGBWNoPublish(1, i, r, g, b, w)
		s.ledBar.SetRGBWNoPublish(2, i, r, g, b, w)
	}
	s.ledBar.Publish()
}

func (s *StreamDeckUI) rotateLEDBarWhite(dialIndex, delta int) {
	if dialIndex > 1 {
		return // Only 2 sections
	}

	section := dialIndex + 1
	avg := 0
	for i := 0; i < 13; i++ {
		val, _ := s.ledBar.GetWhite(section, i)
		avg += val
	}
	avg = avg / 13
	newVal := clamp(avg+delta, 0, 255)

	// Set all white LEDs in section
	for i := 0; i < 13; i++ {
		s.ledBar.SetWhiteNoPublish(section, i, newVal)
	}
	s.ledBar.Publish()
}

func (s *StreamDeckUI) rotateVideoLights(dialIndex, delta int) {
	if dialIndex > 1 {
		return // Only 2 lights
	}

	vl := s.videoLight1
	if dialIndex == 1 {
		vl = s.videoLight2
	}

	on, brightness := vl.GetState()
	newBrightness := clamp(brightness+delta, 0, 100)

	if on {
		vl.TurnOn(newBrightness)
	} else {
		// Update brightness even if off
		vl.TurnOn(newBrightness)
		vl.TurnOff() // Turn back off but brightness is saved
	}
}

// Toggle handlers
func (s *StreamDeckUI) toggleLEDStrip(dialIndex int) {
	// Implementation: toggle between 0 and last value
	// Store in s.lastValues[dialIndex]
}

// ... similar toggle implementations

// clamp restricts value to [min, max]
func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
```

**Validation:**
- Button presses change mode
- Dial rotation updates values
- Dial clicks toggle values
- Touchscreen updates periodically

---

### Phase 7: Integration with main.go

**Goal:** Add Stream Deck mode detection and startup

**File:** `main.go`

```go
// Add import
import "github.com/kevin/office_lights/streamdeck"

// Add after TUI and Web mode detection
useStreamDeck := false
if len(os.Args) > 1 {
	for _, arg := range os.Args[1:] {
		if arg == "streamdeck" {
			useStreamDeck = true
		}
	}
}
if os.Getenv("STREAMDECK") != "" {
	useStreamDeck = true
}

// Add after web server startup
if useStreamDeck {
	go func() {
		log.Println("Starting Stream Deck interface...")
		if err := streamdeck.Run(ledStrip, ledBar, videoLight1, videoLight2); err != nil {
			log.Printf("Stream Deck error: %v", err)
		}
		log.Println("Stream Deck exited")
	}()
}
```

**Validation:**
- `./office_lights streamdeck` starts Stream Deck UI
- Works with other UIs: `./office_lights tui web streamdeck`
- Graceful error if device not found

---

### Phase 8: Polish and Optimization

**Goal:** Improve visuals and performance

**Enhancements:**

1. **Better Button Icons:**
   - Create 120x120 PNG icons for each mode
   - Load from `streamdeck/icons/` directory
   - Overlay with highlight when active

2. **Progress Bars:**
   - Add visual bar graphs to touchscreen sections
   - Show value as percentage of max

3. **Better Fonts:**
   - Use TrueType font instead of basicfont
   - Larger, clearer text

4. **Efficient Rendering:**
   - Only re-render touchscreen if values changed
   - Cache button images
   - Reduce JPEG quality for faster encoding

5. **Error Handling:**
   - Show error messages on touchscreen
   - Handle device disconnect gracefully

**Validation:**
- Icons look professional
- Text is clear and readable
- Performance is smooth (no lag)
- No crashes on errors

---

### Phase 9: Documentation

**Goal:** Document Stream Deck usage

**Update Files:**

1. **CONFIG.md:**
```markdown
### Stream Deck Mode

Run with Stream Deck+ interface:

\`\`\`bash
./office_lights streamdeck
\`\`\`

Requirements:
- Stream Deck+ hardware
- USB connection

Controls:
- Top row buttons: Reserved for future functionality
- Second row buttons: Select mode (LED Strip, LED Bar RGBW, LED Bar White, Video Lights)
- Touchscreen: View current values
- Dials: Rotate to adjust values, click to toggle
\`\`\`

2. **IMPLEMENTATION.md:**
   - Add Phase 13: Stream Deck Interface
   - Document package structure
   - List dependencies

3. **README.md:**
   - Add Stream Deck section (already present)
   - Add usage examples

**Validation:**
- Documentation is clear
- Examples work as documented

---

### Phase 10: Testing

**Goal:** Verify all functionality

**Test Cases:**

1. **Device Detection:**
   - [ ] Device found and opened successfully
   - [ ] Graceful error if device not found
   - [ ] Device resets on startup

2. **Mode Selection:**
   - [ ] Button 4 selects LED Strip mode
   - [ ] Button 5 selects LED Bar RGBW mode
   - [ ] Button 6 selects LED Bar White mode
   - [ ] Button 7 selects Video Lights mode
   - [ ] Active button highlighted

3. **Touchscreen Display:**
   - [ ] LED Strip: Shows R, G, B values
   - [ ] LED Bar RGBW: Shows R, G, B, W values
   - [ ] LED Bar White: Shows Section 1 & 2 averages
   - [ ] Video Lights: Shows Light 1 & 2 brightness and on/off

4. **Dial Rotation:**
   - [ ] Clockwise increases value
   - [ ] Counter-clockwise decreases value
   - [ ] Values clamped to valid range
   - [ ] Changes reflected in touchscreen
   - [ ] MQTT messages published

5. **Dial Click:**
   - [ ] LED Strip: Toggles between 0 and last value
   - [ ] LED Bar: Toggles between 0 and last value
   - [ ] Video Lights: Toggles on/off state

6. **Concurrent Operation:**
   - [ ] Works with TUI simultaneously
   - [ ] Works with Web interface simultaneously
   - [ ] Changes in one UI reflected in Stream Deck
   - [ ] Stream Deck changes reflected in other UIs

7. **Shutdown:**
   - [ ] Ctrl+C shuts down gracefully
   - [ ] Device reset on exit
   - [ ] No resource leaks

**Manual Testing:**
- Test all modes with real hardware
- Verify responsiveness of dials
- Check visual quality of displays

---

## Linux udev Rules

For Linux systems, create `/etc/udev/rules.d/50-streamdeck.rules`:

```
SUBSYSTEM=="usb", ATTRS{idVendor}=="0fd9", ATTRS{idProduct}=="0084", MODE="0666"
```

Then reload rules:
```bash
sudo udevadm control --reload-rules
sudo udevadm trigger
```

---

## Summary

The Stream Deck+ implementation provides:
- 4 operational modes via tactile buttons
- Real-time value display on touchscreen
- Physical dial control with rotation and click
- Concurrent operation with other UIs
- Direct USB HID communication
- Clean integration with existing drivers

After completion, users can control lights via:
1. Terminal (TUI)
2. Web browser
3. Stream Deck+ hardware

All three UIs work simultaneously and share state.

---

## Troubleshooting

**Device not found:**
- Check USB connection
- Verify udev rules (Linux)
- Check user permissions
- Try with sudo (macOS/Linux)

**Rendering issues:**
- Check image dimensions (button: 120x120, touch: 800x100)
- Verify JPEG encoding quality
- Test with simple solid colors first

**Dial sensitivity:**
- Adjust multiplier (currently × 5)
- Add acceleration for faster changes
- Implement dead zone if needed

**Performance:**
- Reduce touchscreen update frequency
- Optimize image rendering
- Profile with pprof

---

## Next Steps

After implementation:
1. Test with real hardware
2. Gather user feedback
3. Consider enhancements:
   - Scene presets on bottom row buttons
   - Touchscreen tap support
   - Custom dial sensitivity settings
   - Animation effects
