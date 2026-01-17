package ledstrip

import (
	"encoding/json"
	"testing"

	"github.com/kevin/office_lights/mqtt"
)

func TestNewLEDStrip(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	strip := NewLEDStrip(mock, "test/topic")

	if strip == nil {
		t.Fatal("NewLEDStrip returned nil")
	}

	r, g, b := strip.GetColor()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("Expected initial color (0,0,0), got (%d,%d,%d)", r, g, b)
	}
}

func TestSetColor(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	strip := NewLEDStrip(mock, "test/topic")

	tests := []struct {
		name      string
		r, g, b   int
		wantError bool
	}{
		{"Valid color", 100, 150, 200, false},
		{"Max values", 255, 255, 255, false},
		{"Min values", 0, 0, 0, false},
		{"Red too high", 256, 0, 0, true},
		{"Red negative", -1, 0, 0, true},
		{"Green too high", 0, 256, 0, true},
		{"Blue negative", 0, 0, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.Clear()
			err := strip.SetColor(tt.r, tt.g, tt.b)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error for SetColor(%d,%d,%d), got nil", tt.r, tt.g, tt.b)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Check color was set
				r, g, b := strip.GetColor()
				if r != tt.r || g != tt.g || b != tt.b {
					t.Errorf("Expected color (%d,%d,%d), got (%d,%d,%d)", tt.r, tt.g, tt.b, r, g, b)
				}

				// Check message was published
				if mock.MessageCount() != 1 {
					t.Errorf("Expected 1 message, got %d", mock.MessageCount())
				}
			}
		})
	}
}

func TestMessageFormat(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	strip := NewLEDStrip(mock, "test/topic")

	err := strip.SetColor(100, 150, 200)
	if err != nil {
		t.Fatalf("SetColor failed: %v", err)
	}

	msg := mock.GetLastMessage()
	if msg == nil {
		t.Fatal("No message published")
	}

	if msg.Topic != "test/topic" {
		t.Errorf("Expected topic 'test/topic', got '%s'", msg.Topic)
	}

	// Parse JSON payload
	var result sequenceMessage
	if err := json.Unmarshal(msg.Payload.([]byte), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result.Sequence != "fill" {
		t.Errorf("Expected sequence 'fill', got '%s'", result.Sequence)
	}

	if result.Data.R != 100 || result.Data.G != 150 || result.Data.B != 200 {
		t.Errorf("Expected RGB (100,150,200), got (%d,%d,%d)",
			result.Data.R, result.Data.G, result.Data.B)
	}
}

func TestTurnOff(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	strip := NewLEDStrip(mock, "test/topic")

	// First set a color
	strip.SetColor(100, 150, 200)
	mock.Clear()

	// Turn off
	err := strip.TurnOff()
	if err != nil {
		t.Fatalf("TurnOff failed: %v", err)
	}

	r, g, b := strip.GetColor()
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("Expected color (0,0,0) after TurnOff, got (%d,%d,%d)", r, g, b)
	}

	if mock.MessageCount() != 1 {
		t.Errorf("Expected 1 message, got %d", mock.MessageCount())
	}
}

func TestSetBrightness(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	strip := NewLEDStrip(mock, "test/topic")

	// Set initial color
	strip.SetColor(200, 100, 50)
	mock.Clear()

	// Set brightness to 50%
	err := strip.SetBrightness(50)
	if err != nil {
		t.Fatalf("SetBrightness failed: %v", err)
	}

	r, g, b := strip.GetColor()
	if r != 100 || g != 50 || b != 25 {
		t.Errorf("Expected color (100,50,25) at 50%% brightness, got (%d,%d,%d)", r, g, b)
	}

	// Test invalid brightness
	err = strip.SetBrightness(101)
	if err == nil {
		t.Error("Expected error for brightness > 100")
	}

	err = strip.SetBrightness(-1)
	if err == nil {
		t.Error("Expected error for brightness < 0")
	}
}

func TestPresetColors(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	strip := NewLEDStrip(mock, "test/topic")

	tests := []struct {
		name           string
		setFunc        func() error
		expectedR, g, b int
	}{
		{"Red", strip.SetRed, 255, 0, 0},
		{"Green", strip.SetGreen, 0, 255, 0},
		{"Blue", strip.SetBlue, 0, 0, 255},
		{"White", strip.SetWhite, 255, 255, 255},
		{"Yellow", strip.SetYellow, 255, 255, 0},
		{"Cyan", strip.SetCyan, 0, 255, 255},
		{"Magenta", strip.SetMagenta, 255, 0, 255},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.Clear()
			err := tt.setFunc()
			if err != nil {
				t.Fatalf("Preset color failed: %v", err)
			}

			r, g, b := strip.GetColor()
			if r != tt.expectedR || g != tt.g || b != tt.b {
				t.Errorf("Expected color (%d,%d,%d), got (%d,%d,%d)",
					tt.expectedR, tt.g, tt.b, r, g, b)
			}

			if mock.MessageCount() != 1 {
				t.Errorf("Expected 1 message, got %d", mock.MessageCount())
			}
		})
	}
}
