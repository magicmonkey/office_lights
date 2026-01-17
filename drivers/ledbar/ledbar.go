package ledbar

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Publisher defines the interface for publishing MQTT messages
type Publisher interface {
	Publish(topic string, payload interface{}) error
}

// StateStore defines the interface for persistent state storage
type StateStore interface {
	SaveLEDBarChannels(ledbarID int, channels []int) error
}

// LEDBar represents an RGBW LED bar controller
// The bar has:
// - 6 RGBW LEDs in section 1
// - 13 white LEDs in section 1
// - 6 RGBW LEDs in section 2
// - 13 white LEDs in section 2
type LEDBar struct {
	rgbw1     [6][4]int // First set of 6 RGBW LEDs (R, G, B, W)
	white1    [13]int   // First set of 13 white LEDs
	rgbw2     [6][4]int // Second set of 6 RGBW LEDs (R, G, B, W)
	white2    [13]int   // Second set of 13 white LEDs
	barID     int
	publisher Publisher
	topic     string
	store     StateStore
}

// NewLEDBar creates a new LED bar controller with default state (all off)
func NewLEDBar(barID int, publisher Publisher, topic string) (*LEDBar, error) {
	channels := make([]int, 77)
	return NewLEDBarWithState(barID, publisher, topic, nil, channels)
}

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

// SetRGBW sets the RGBW values for a specific LED in a section
// section: 1 or 2
// index: 0-5 (which RGBW LED)
// r, g, b, w: 0-255
func (l *LEDBar) SetRGBW(section int, index int, r, g, b, w int) error {
	if section != 1 && section != 2 {
		return fmt.Errorf("section must be 1 or 2, got %d", section)
	}
	if index < 0 || index > 5 {
		return fmt.Errorf("index must be between 0 and 5, got %d", index)
	}
	if err := validateValue(r); err != nil {
		return fmt.Errorf("red: %w", err)
	}
	if err := validateValue(g); err != nil {
		return fmt.Errorf("green: %w", err)
	}
	if err := validateValue(b); err != nil {
		return fmt.Errorf("blue: %w", err)
	}
	if err := validateValue(w); err != nil {
		return fmt.Errorf("white: %w", err)
	}

	if section == 1 {
		l.rgbw1[index][0] = r
		l.rgbw1[index][1] = g
		l.rgbw1[index][2] = b
		l.rgbw1[index][3] = w
	} else {
		l.rgbw2[index][0] = r
		l.rgbw2[index][1] = g
		l.rgbw2[index][2] = b
		l.rgbw2[index][3] = w
	}

	return l.Publish()
}

// SetWhite sets the white LED value for a specific LED in a section
// section: 1 or 2
// index: 0-12 (which white LED)
// value: 0-255
func (l *LEDBar) SetWhite(section int, index int, value int) error {
	if section != 1 && section != 2 {
		return fmt.Errorf("section must be 1 or 2, got %d", section)
	}
	if index < 0 || index > 12 {
		return fmt.Errorf("index must be between 0 and 12, got %d", index)
	}
	if err := validateValue(value); err != nil {
		return fmt.Errorf("value: %w", err)
	}

	if section == 1 {
		l.white1[index] = value
	} else {
		l.white2[index] = value
	}

	return l.Publish()
}

// GetRGBW returns the RGBW values for a specific LED in a section
func (l *LEDBar) GetRGBW(section int, index int) (int, int, int, int, error) {
	if section != 1 && section != 2 {
		return 0, 0, 0, 0, fmt.Errorf("section must be 1 or 2, got %d", section)
	}
	if index < 0 || index > 5 {
		return 0, 0, 0, 0, fmt.Errorf("index must be between 0 and 5, got %d", index)
	}

	if section == 1 {
		return l.rgbw1[index][0], l.rgbw1[index][1], l.rgbw1[index][2], l.rgbw1[index][3], nil
	}
	return l.rgbw2[index][0], l.rgbw2[index][1], l.rgbw2[index][2], l.rgbw2[index][3], nil
}

// GetWhite returns the white LED value for a specific LED in a section
func (l *LEDBar) GetWhite(section int, index int) (int, error) {
	if section != 1 && section != 2 {
		return 0, fmt.Errorf("section must be 1 or 2, got %d", section)
	}
	if index < 0 || index > 12 {
		return 0, fmt.Errorf("index must be between 0 and 12, got %d", index)
	}

	if section == 1 {
		return l.white1[index], nil
	}
	return l.white2[index], nil
}

// TurnOffSection turns off all LEDs in a section
func (l *LEDBar) TurnOffSection(section int) error {
	if section != 1 && section != 2 {
		return fmt.Errorf("section must be 1 or 2, got %d", section)
	}

	if section == 1 {
		for i := range l.rgbw1 {
			for j := range l.rgbw1[i] {
				l.rgbw1[i][j] = 0
			}
		}
		for i := range l.white1 {
			l.white1[i] = 0
		}
	} else {
		for i := range l.rgbw2 {
			for j := range l.rgbw2[i] {
				l.rgbw2[i][j] = 0
			}
		}
		for i := range l.white2 {
			l.white2[i] = 0
		}
	}

	return l.Publish()
}

// TurnOffAll turns off all LEDs on the bar
func (l *LEDBar) TurnOffAll() error {
	for i := range l.rgbw1 {
		for j := range l.rgbw1[i] {
			l.rgbw1[i][j] = 0
		}
	}
	for i := range l.white1 {
		l.white1[i] = 0
	}
	for i := range l.rgbw2 {
		for j := range l.rgbw2[i] {
			l.rgbw2[i][j] = 0
		}
	}
	for i := range l.white2 {
		l.white2[i] = 0
	}

	return l.Publish()
}

// SetAllRGBW sets all RGBW LEDs to the same color
func (l *LEDBar) SetAllRGBW(r, g, b, w int) error {
	if err := validateValue(r); err != nil {
		return fmt.Errorf("red: %w", err)
	}
	if err := validateValue(g); err != nil {
		return fmt.Errorf("green: %w", err)
	}
	if err := validateValue(b); err != nil {
		return fmt.Errorf("blue: %w", err)
	}
	if err := validateValue(w); err != nil {
		return fmt.Errorf("white: %w", err)
	}

	for i := range l.rgbw1 {
		l.rgbw1[i][0] = r
		l.rgbw1[i][1] = g
		l.rgbw1[i][2] = b
		l.rgbw1[i][3] = w
	}
	for i := range l.rgbw2 {
		l.rgbw2[i][0] = r
		l.rgbw2[i][1] = g
		l.rgbw2[i][2] = b
		l.rgbw2[i][3] = w
	}

	return l.Publish()
}

// SetAllWhite sets all white LEDs to the same value
func (l *LEDBar) SetAllWhite(value int) error {
	if err := validateValue(value); err != nil {
		return fmt.Errorf("value: %w", err)
	}

	for i := range l.white1 {
		l.white1[i] = value
	}
	for i := range l.white2 {
		l.white2[i] = value
	}

	return l.Publish()
}

// Publish formats and publishes the current state to MQTT
func (l *LEDBar) Publish() error {
	payload := l.formatMessage()

	if err := l.publisher.Publish(l.topic, payload); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	// Save state to storage after successful publish
	if l.store != nil {
		channels := l.getChannels()
		if err := l.store.SaveLEDBarChannels(l.barID, channels); err != nil {
			// Log error but don't fail the operation
			log.Printf("Warning: Failed to save LED bar state: %v", err)
		}
	}

	return nil
}

// formatMessage creates the comma-separated message for the LED bar
// Message structure (77 values total):
// - Values 0-23: 6 RGBW LEDs (4 values each: R,G,B,W)
// - Values 24-36: 13 white LEDs (1 value each)
// - Values 37-39: 3 ignored values (set to 0)
// - Values 40-63: 6 RGBW LEDs (4 values each: R,G,B,W)
// - Values 64-76: 13 white LEDs (1 value each)
func (l *LEDBar) formatMessage() string {
	values := make([]string, 0, 77)

	// First section: 6 RGBW LEDs (24 values)
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			values = append(values, strconv.Itoa(l.rgbw1[i][j]))
		}
	}

	// First section: 13 white LEDs (13 values)
	for i := 0; i < 13; i++ {
		values = append(values, strconv.Itoa(l.white1[i]))
	}

	// 3 ignored values
	values = append(values, "0", "0", "0")

	// Second section: 6 RGBW LEDs (24 values)
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			values = append(values, strconv.Itoa(l.rgbw2[i][j]))
		}
	}

	// Second section: 13 white LEDs (13 values)
	for i := 0; i < 13; i++ {
		values = append(values, strconv.Itoa(l.white2[i]))
	}

	return strings.Join(values, ",")
}

// validateValue validates that a value is in the valid range (0-255)
func validateValue(value int) error {
	if value < 0 || value > 255 {
		return fmt.Errorf("value must be between 0 and 255, got %d", value)
	}
	return nil
}

// GetBarID returns the bar ID
func (l *LEDBar) GetBarID() int {
	return l.barID
}

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
