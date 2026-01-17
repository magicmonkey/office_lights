package videolight

import (
	"fmt"
	"strings"
)

// Publisher defines the interface for publishing MQTT messages
type Publisher interface {
	Publish(topic string, payload interface{}) error
}

// VideoLight represents a video light controller
type VideoLight struct {
	on         bool
	brightness int
	lightID    int
	publisher  Publisher
	topic      string
}

// NewVideoLight creates a new video light controller
func NewVideoLight(lightID int, publisher Publisher, topic string) (*VideoLight, error) {
	if lightID < 1 {
		return nil, fmt.Errorf("lightID must be positive, got %d", lightID)
	}

	return &VideoLight{
		on:         false,
		brightness: 0,
		lightID:    lightID,
		publisher:  publisher,
		topic:      topic,
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
