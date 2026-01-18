package videolight

import (
	"testing"

	"github.com/kevin/office_lights/mqtt"
)

func TestNewVideoLight(t *testing.T) {
	mock := mqtt.NewMockPublisher()

	tests := []struct {
		name      string
		lightID   int
		wantError bool
	}{
		{"Valid ID 1", 1, false},
		{"Valid ID 2", 2, false},
		{"Valid ID 100", 100, false},
		{"Invalid ID 0", 0, true},
		{"Invalid ID negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			light, err := NewVideoLight(tt.lightID, mock, "test/topic")

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error for lightID %d, got nil", tt.lightID)
				}
				if light != nil {
					t.Errorf("Expected nil light for invalid ID, got %+v", light)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if light == nil {
					t.Fatal("Expected non-nil light for valid ID")
				}

				// Check initial state
				on, brightness := light.GetState()
				if on {
					t.Error("Expected light to be off initially")
				}
				if brightness != 0 {
					t.Errorf("Expected initial brightness 0, got %d", brightness)
				}

				if light.GetLightID() != tt.lightID {
					t.Errorf("Expected lightID %d, got %d", tt.lightID, light.GetLightID())
				}
			}
		})
	}
}

func TestSetState(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	light, _ := NewVideoLight(1, mock, "test/topic")

	tests := []struct {
		name       string
		on         bool
		brightness int
		wantError  bool
	}{
		{"Turn on at 50%", true, 50, false},
		{"Turn on at 100%", true, 100, false},
		{"Turn on at 0%", true, 0, false},
		{"Turn off", false, 0, false},
		{"Brightness too high", true, 101, true},
		{"Brightness negative", true, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.Clear()
			err := light.SetState(tt.on, tt.brightness)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error for brightness %d, got nil", tt.brightness)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check state was set
				on, brightness := light.GetState()
				if on != tt.on {
					t.Errorf("Expected on=%v, got %v", tt.on, on)
				}
				if brightness != tt.brightness {
					t.Errorf("Expected brightness %d, got %d", tt.brightness, brightness)
				}

				// Check message was published
				if mock.MessageCount() != 1 {
					t.Errorf("Expected 1 message, got %d", mock.MessageCount())
				}
			}
		})
	}
}

func TestTurnOn(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	light, _ := NewVideoLight(1, mock, "test/topic")

	err := light.TurnOn(75)
	if err != nil {
		t.Fatalf("TurnOn failed: %v", err)
	}

	on, brightness := light.GetState()
	if !on {
		t.Error("Expected light to be on")
	}
	if brightness != 75 {
		t.Errorf("Expected brightness 75, got %d", brightness)
	}

	if mock.MessageCount() != 1 {
		t.Errorf("Expected 1 message, got %d", mock.MessageCount())
	}
}

func TestTurnOff(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	light, _ := NewVideoLight(1, mock, "test/topic")

	// First turn on
	light.TurnOn(50)
	mock.Clear()

	// Then turn off
	err := light.TurnOff()
	if err != nil {
		t.Fatalf("TurnOff failed: %v", err)
	}

	on, brightness := light.GetState()
	if on {
		t.Error("Expected light to be off")
	}
	if brightness != 50 {
		t.Errorf("Expected brightness to be preserved at 50, got %d", brightness)
	}

	if mock.MessageCount() != 1 {
		t.Errorf("Expected 1 message, got %d", mock.MessageCount())
	}
}

func TestSetBrightness(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	light, _ := NewVideoLight(1, mock, "test/topic")

	// Set initial state to on
	light.SetState(true, 50)
	mock.Clear()

	// Change brightness
	err := light.SetBrightness(75)
	if err != nil {
		t.Fatalf("SetBrightness failed: %v", err)
	}

	on, brightness := light.GetState()
	if !on {
		t.Error("Expected light to remain on")
	}
	if brightness != 75 {
		t.Errorf("Expected brightness 75, got %d", brightness)
	}
}

func TestMessageFormat(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	light, _ := NewVideoLight(1, mock, "test/topic")

	tests := []struct {
		name           string
		on             bool
		brightness     int
		expectedPayload string
	}{
		{"On at 50%", true, 50, "set,true,50"},
		{"On at 100%", true, 100, "set,true,100"},
		{"On at 0%", true, 0, "set,true,0"},
		{"Off", false, 0, "set,false,0"},
		{"Off with brightness", false, 50, "set,false,50"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.Clear()
			light.SetState(tt.on, tt.brightness)

			msg := mock.GetLastMessage()
			if msg == nil {
				t.Fatal("No message published")
			}

			if msg.Topic != "test/topic" {
				t.Errorf("Expected topic 'test/topic', got '%s'", msg.Topic)
			}

			payload := msg.Payload.(string)
			if payload != tt.expectedPayload {
				t.Errorf("Expected payload '%s', got '%s'", tt.expectedPayload, payload)
			}
		})
	}
}

func TestMultipleLights(t *testing.T) {
	mock := mqtt.NewMockPublisher()

	light1, _ := NewVideoLight(1, mock, "topic1")
	light2, _ := NewVideoLight(2, mock, "topic2")

	light1.TurnOn(50)
	light2.TurnOn(75)

	if light1.GetLightID() != 1 {
		t.Errorf("Light 1 has wrong ID: %d", light1.GetLightID())
	}
	if light2.GetLightID() != 2 {
		t.Errorf("Light 2 has wrong ID: %d", light2.GetLightID())
	}

	on1, brightness1 := light1.GetState()
	on2, brightness2 := light2.GetState()

	if !on1 || brightness1 != 50 {
		t.Errorf("Light 1 state incorrect: on=%v, brightness=%d", on1, brightness1)
	}
	if !on2 || brightness2 != 75 {
		t.Errorf("Light 2 state incorrect: on=%v, brightness=%d", on2, brightness2)
	}

	// Should have 2 messages (one for each light)
	if mock.MessageCount() != 2 {
		t.Errorf("Expected 2 messages, got %d", mock.MessageCount())
	}
}
