package ledstrip

import (
	"encoding/json"
	"fmt"
	"log"
)

// Publisher defines the interface for publishing MQTT messages
type Publisher interface {
	Publish(topic string, payload interface{}) error
}

// StateStore defines the interface for persistent state storage
type StateStore interface {
	SaveLEDStripState(id int, r, g, b int) error
}

// LEDStrip represents an RGB LED strip controller
type LEDStrip struct {
	r         int
	g         int
	b         int
	publisher Publisher
	topic     string
	store     StateStore
	id        int
}

// sequenceMessage represents the JSON structure for LED strip commands
type sequenceMessage struct {
	Sequence string      `json:"sequence"`
	Data     sequenceData `json:"data"`
}

// sequenceData represents the RGB data in the sequence message
type sequenceData struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

// NewLEDStrip creates a new LED strip controller with default state (all off)
func NewLEDStrip(publisher Publisher, topic string) *LEDStrip {
	return NewLEDStripWithState(publisher, topic, nil, 0, 0, 0, 0)
}

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

// SetColor sets the RGB color values and publishes the update
func (l *LEDStrip) SetColor(r, g, b int) error {
	if err := validateRGB(r, g, b); err != nil {
		return err
	}

	l.r = r
	l.g = g
	l.b = b

	return l.Publish()
}

// GetColor returns the current RGB color values
func (l *LEDStrip) GetColor() (int, int, int) {
	return l.r, l.g, l.b
}

// TurnOff turns off the LED strip by setting all colors to 0
func (l *LEDStrip) TurnOff() error {
	return l.SetColor(0, 0, 0)
}

// SetBrightness adjusts the current color by a percentage (0-100)
func (l *LEDStrip) SetBrightness(percentage int) error {
	if percentage < 0 || percentage > 100 {
		return fmt.Errorf("brightness must be between 0 and 100, got %d", percentage)
	}

	scale := float64(percentage) / 100.0
	r := int(float64(l.r) * scale)
	g := int(float64(l.g) * scale)
	b := int(float64(l.b) * scale)

	return l.SetColor(r, g, b)
}

// Publish formats and publishes the current state to MQTT
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

// formatMessage creates the JSON message for the LED strip
func (l *LEDStrip) formatMessage() ([]byte, error) {
	msg := sequenceMessage{
		Sequence: "fill",
		Data: sequenceData{
			R: l.r,
			G: l.g,
			B: l.b,
		},
	}

	return json.Marshal(msg)
}

// validateRGB validates that RGB values are in the valid range (0-255)
func validateRGB(r, g, b int) error {
	if r < 0 || r > 255 {
		return fmt.Errorf("red value must be between 0 and 255, got %d", r)
	}
	if g < 0 || g > 255 {
		return fmt.Errorf("green value must be between 0 and 255, got %d", g)
	}
	if b < 0 || b > 255 {
		return fmt.Errorf("blue value must be between 0 and 255, got %d", b)
	}
	return nil
}

// Preset color methods for convenience

// SetRed sets the strip to red
func (l *LEDStrip) SetRed() error {
	return l.SetColor(255, 0, 0)
}

// SetGreen sets the strip to green
func (l *LEDStrip) SetGreen() error {
	return l.SetColor(0, 255, 0)
}

// SetBlue sets the strip to blue
func (l *LEDStrip) SetBlue() error {
	return l.SetColor(0, 0, 255)
}

// SetWhite sets the strip to white
func (l *LEDStrip) SetWhite() error {
	return l.SetColor(255, 255, 255)
}

// SetYellow sets the strip to yellow
func (l *LEDStrip) SetYellow() error {
	return l.SetColor(255, 255, 0)
}

// SetCyan sets the strip to cyan
func (l *LEDStrip) SetCyan() error {
	return l.SetColor(0, 255, 255)
}

// SetMagenta sets the strip to magenta
func (l *LEDStrip) SetMagenta() error {
	return l.SetColor(255, 0, 255)
}
