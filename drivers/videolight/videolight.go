package videolight

import (
	"fmt"
	"log"
	"strings"
)

// Publisher defines the interface for publishing MQTT messages
type Publisher interface {
	Publish(topic string, payload interface{}) error
}

// StateStore defines the interface for persistent state storage
type StateStore interface {
	SaveVideoLightState(id int, on bool, brightness int) error
}

// VideoLight represents a video light controller
type VideoLight struct {
	on         bool
	brightness int
	lightID    int
	publisher  Publisher
	topic      string
	store      StateStore
}

// NewVideoLight creates a new video light controller with default state (off)
func NewVideoLight(lightID int, publisher Publisher, topic string) (*VideoLight, error) {
	return NewVideoLightWithState(lightID, publisher, topic, nil, false, 0)
}

// NewVideoLightWithState creates video light with initial state from storage
func NewVideoLightWithState(lightID int, publisher Publisher, topic string, store StateStore, on bool, brightness int) (*VideoLight, error) {
	if lightID < 1 {
		return nil, fmt.Errorf("lightID must be positive, got %d", lightID)
	}

	// Validate and fix any invalid stored state
	if err := validateBrightness(brightness); err != nil {
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

// SetState sets the on/off state and brightness, then publishes the update
func (v *VideoLight) SetState(on bool, brightness int) error {
	if err := validateBrightness(brightness); err != nil {
		return err
	}

	v.on = on
	v.brightness = brightness

	return v.Publish()
}

// TurnOn turns on the light at the specified brightness
func (v *VideoLight) TurnOn(brightness int) error {
	return v.SetState(true, brightness)
}

// TurnOff turns off the light
func (v *VideoLight) TurnOff() error {
	return v.SetState(false, 0)
}

// SetBrightness sets the brightness while maintaining the current on/off state
func (v *VideoLight) SetBrightness(brightness int) error {
	return v.SetState(v.on, brightness)
}

// GetState returns the current on/off state and brightness
func (v *VideoLight) GetState() (bool, int) {
	return v.on, v.brightness
}

// IsOn returns whether the light is currently on
func (v *VideoLight) IsOn() bool {
	return v.on
}

// Brightness returns the current brightness value
func (v *VideoLight) Brightness() int {
	return v.brightness
}

// GetLightID returns the light ID
func (v *VideoLight) GetLightID() int {
	return v.lightID
}

// Publish formats and publishes the current state to MQTT
func (v *VideoLight) Publish() error {
	payload := v.formatMessage()

	if err := v.publisher.Publish(v.topic, payload); err != nil {
		return fmt.Errorf("failed to publish: %w", err)
	}

	// Save state to storage after successful publish
	if v.store != nil {
		// Convert driver ID (1, 2) to database ID (0, 1)
		dbID := v.lightID - 1
		if err := v.store.SaveVideoLightState(dbID, v.on, v.brightness); err != nil {
			// Log error but don't fail the operation
			log.Printf("Warning: Failed to save video light state: %v", err)
		}
	}

	return nil
}

// formatMessage creates the message string for the video light
// Format: set,<on>,<brightness>
// Example: set,true,50
func (v *VideoLight) formatMessage() string {
	var builder strings.Builder
	builder.WriteString("set,")

	if v.on {
		builder.WriteString("true")
	} else {
		builder.WriteString("false")
	}

	builder.WriteString(",")
	builder.WriteString(fmt.Sprintf("%d", v.brightness))

	return builder.String()
}

// validateBrightness validates that brightness is in the valid range (0-100)
func validateBrightness(brightness int) error {
	if brightness < 0 || brightness > 100 {
		return fmt.Errorf("brightness must be between 0 and 100, got %d", brightness)
	}
	return nil
}
