package ledbar

import (
	"strings"
	"testing"

	"github.com/kevin/office_lights/mqtt"
)

func TestNewLEDBar(t *testing.T) {
	mock := mqtt.NewMockPublisher()

	tests := []struct {
		name      string
		barID     int
		wantError bool
	}{
		{"Valid ID 0", 0, false},
		{"Valid ID 1", 1, false},
		{"Valid ID 100", 100, false},
		{"Invalid ID negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bar, err := NewLEDBar(tt.barID, mock, "test/topic")

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error for barID %d, got nil", tt.barID)
				}
				if bar != nil {
					t.Errorf("Expected nil bar for invalid ID, got %+v", bar)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if bar == nil {
					t.Fatal("Expected non-nil bar for valid ID")
				}

				if bar.GetBarID() != tt.barID {
					t.Errorf("Expected barID %d, got %d", tt.barID, bar.GetBarID())
				}

				// Verify all LEDs are initialized to 0
				for i := 0; i < 6; i++ {
					r, g, b, w, _ := bar.GetRGBW(1, i)
					if r != 0 || g != 0 || b != 0 || w != 0 {
						t.Errorf("Section 1 RGBW[%d] not initialized to 0, got (%d,%d,%d,%d)", i, r, g, b, w)
					}

					r, g, b, w, _ = bar.GetRGBW(2, i)
					if r != 0 || g != 0 || b != 0 || w != 0 {
						t.Errorf("Section 2 RGBW[%d] not initialized to 0, got (%d,%d,%d,%d)", i, r, g, b, w)
					}
				}

				for i := 0; i < 13; i++ {
					val, _ := bar.GetWhite(1, i)
					if val != 0 {
						t.Errorf("Section 1 White[%d] not initialized to 0, got %d", i, val)
					}

					val, _ = bar.GetWhite(2, i)
					if val != 0 {
						t.Errorf("Section 2 White[%d] not initialized to 0, got %d", i, val)
					}
				}
			}
		})
	}
}

func TestSetRGBW(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	tests := []struct {
		name      string
		section   int
		index     int
		r, g, b, w int
		wantError bool
	}{
		{"Valid section 1", 1, 0, 100, 150, 200, 50, false},
		{"Valid section 2", 2, 5, 255, 0, 128, 200, false},
		{"Max values", 1, 3, 255, 255, 255, 255, false},
		{"Min values", 2, 0, 0, 0, 0, 0, false},
		{"Invalid section 0", 0, 0, 100, 100, 100, 100, true},
		{"Invalid section 3", 3, 0, 100, 100, 100, 100, true},
		{"Invalid index -1", 1, -1, 100, 100, 100, 100, true},
		{"Invalid index 6", 1, 6, 100, 100, 100, 100, true},
		{"R too high", 1, 0, 256, 100, 100, 100, true},
		{"G negative", 1, 0, 100, -1, 100, 100, true},
		{"B too high", 1, 0, 100, 100, 256, 100, true},
		{"W negative", 1, 0, 100, 100, 100, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.Clear()
			err := bar.SetRGBW(tt.section, tt.index, tt.r, tt.g, tt.b, tt.w)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify value was set
				r, g, b, w, err := bar.GetRGBW(tt.section, tt.index)
				if err != nil {
					t.Fatalf("GetRGBW failed: %v", err)
				}
				if r != tt.r || g != tt.g || b != tt.b || w != tt.w {
					t.Errorf("Expected RGBW (%d,%d,%d,%d), got (%d,%d,%d,%d)",
						tt.r, tt.g, tt.b, tt.w, r, g, b, w)
				}

				// Check message was published
				if mock.MessageCount() != 1 {
					t.Errorf("Expected 1 message, got %d", mock.MessageCount())
				}
			}
		})
	}
}

func TestSetWhite(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	tests := []struct {
		name      string
		section   int
		index     int
		value     int
		wantError bool
	}{
		{"Valid section 1", 1, 0, 100, false},
		{"Valid section 2", 2, 12, 255, false},
		{"Min value", 1, 5, 0, false},
		{"Invalid section", 0, 0, 100, true},
		{"Invalid index -1", 1, -1, 100, true},
		{"Invalid index 13", 1, 13, 100, true},
		{"Value too high", 1, 0, 256, true},
		{"Value negative", 1, 0, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.Clear()
			err := bar.SetWhite(tt.section, tt.index, tt.value)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Verify value was set
				val, err := bar.GetWhite(tt.section, tt.index)
				if err != nil {
					t.Fatalf("GetWhite failed: %v", err)
				}
				if val != tt.value {
					t.Errorf("Expected value %d, got %d", tt.value, val)
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
	bar, _ := NewLEDBar(0, mock, "test/topic")

	// Set some values
	bar.SetRGBW(1, 0, 10, 20, 30, 40)
	bar.SetWhite(1, 0, 50)
	bar.SetRGBW(2, 0, 60, 70, 80, 90)
	bar.SetWhite(2, 0, 100)

	msg := mock.GetLastMessage()
	if msg == nil {
		t.Fatal("No message published")
	}

	if msg.Topic != "test/topic" {
		t.Errorf("Expected topic 'test/topic', got '%s'", msg.Topic)
	}

	payload := msg.Payload.(string)
	values := strings.Split(payload, ",")

	// Should have exactly 77 values
	if len(values) != 77 {
		t.Fatalf("Expected 77 values, got %d", len(values))
	}

	// Check first RGBW LED (section 1, index 0)
	// Position 0-3: R,G,B,W = 10,20,30,40
	if values[0] != "10" || values[1] != "20" || values[2] != "30" || values[3] != "40" {
		t.Errorf("First RGBW LED incorrect: got %s,%s,%s,%s", values[0], values[1], values[2], values[3])
	}

	// Check first white LED (section 1, index 0)
	// Position 24: value = 50
	if values[24] != "50" {
		t.Errorf("First white LED incorrect: expected 50, got %s", values[24])
	}

	// Check ignored values (positions 37-39)
	if values[37] != "0" || values[38] != "0" || values[39] != "0" {
		t.Errorf("Ignored values incorrect: got %s,%s,%s", values[37], values[38], values[39])
	}

	// Check first RGBW LED in section 2 (index 0)
	// Position 40-43: R,G,B,W = 60,70,80,90
	if values[40] != "60" || values[41] != "70" || values[42] != "80" || values[43] != "90" {
		t.Errorf("Section 2 first RGBW LED incorrect: got %s,%s,%s,%s", values[40], values[41], values[42], values[43])
	}

	// Check first white LED in section 2 (index 0)
	// Position 64: value = 100
	if values[64] != "100" {
		t.Errorf("Section 2 first white LED incorrect: expected 100, got %s", values[64])
	}
}

func TestTurnOffSection(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	// Set some values in section 1
	bar.SetRGBW(1, 0, 100, 100, 100, 100)
	bar.SetWhite(1, 0, 100)

	// Set some values in section 2
	bar.SetRGBW(2, 0, 200, 200, 200, 200)
	bar.SetWhite(2, 0, 200)

	mock.Clear()

	// Turn off section 1
	err := bar.TurnOffSection(1)
	if err != nil {
		t.Fatalf("TurnOffSection failed: %v", err)
	}

	// Verify section 1 is off
	for i := 0; i < 6; i++ {
		r, g, b, w, _ := bar.GetRGBW(1, i)
		if r != 0 || g != 0 || b != 0 || w != 0 {
			t.Errorf("Section 1 RGBW[%d] not off, got (%d,%d,%d,%d)", i, r, g, b, w)
		}
	}
	for i := 0; i < 13; i++ {
		val, _ := bar.GetWhite(1, i)
		if val != 0 {
			t.Errorf("Section 1 White[%d] not off, got %d", i, val)
		}
	}

	// Verify section 2 is still on
	r, g, b, w, _ := bar.GetRGBW(2, 0)
	if r != 200 || g != 200 || b != 200 || w != 200 {
		t.Errorf("Section 2 RGBW changed unexpectedly")
	}
	val, _ := bar.GetWhite(2, 0)
	if val != 200 {
		t.Errorf("Section 2 white changed unexpectedly")
	}

	// Test invalid section
	err = bar.TurnOffSection(0)
	if err == nil {
		t.Error("Expected error for invalid section")
	}
}

func TestTurnOffAll(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	// Set values in both sections
	bar.SetRGBW(1, 0, 100, 100, 100, 100)
	bar.SetWhite(1, 0, 100)
	bar.SetRGBW(2, 0, 200, 200, 200, 200)
	bar.SetWhite(2, 0, 200)

	mock.Clear()

	err := bar.TurnOffAll()
	if err != nil {
		t.Fatalf("TurnOffAll failed: %v", err)
	}

	// Verify everything is off
	for section := 1; section <= 2; section++ {
		for i := 0; i < 6; i++ {
			r, g, b, w, _ := bar.GetRGBW(section, i)
			if r != 0 || g != 0 || b != 0 || w != 0 {
				t.Errorf("Section %d RGBW[%d] not off, got (%d,%d,%d,%d)", section, i, r, g, b, w)
			}
		}
		for i := 0; i < 13; i++ {
			val, _ := bar.GetWhite(section, i)
			if val != 0 {
				t.Errorf("Section %d White[%d] not off, got %d", section, i, val)
			}
		}
	}
}

func TestSetAllRGBW(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	err := bar.SetAllRGBW(100, 150, 200, 50)
	if err != nil {
		t.Fatalf("SetAllRGBW failed: %v", err)
	}

	// Verify all RGBW LEDs have the same value
	for section := 1; section <= 2; section++ {
		for i := 0; i < 6; i++ {
			r, g, b, w, _ := bar.GetRGBW(section, i)
			if r != 100 || g != 150 || b != 200 || w != 50 {
				t.Errorf("Section %d RGBW[%d] incorrect, got (%d,%d,%d,%d)", section, i, r, g, b, w)
			}
		}
	}

	// Test invalid values
	err = bar.SetAllRGBW(256, 0, 0, 0)
	if err == nil {
		t.Error("Expected error for invalid R value")
	}
}

func TestSetAllWhite(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	err := bar.SetAllWhite(150)
	if err != nil {
		t.Fatalf("SetAllWhite failed: %v", err)
	}

	// Verify all white LEDs have the same value
	for section := 1; section <= 2; section++ {
		for i := 0; i < 13; i++ {
			val, _ := bar.GetWhite(section, i)
			if val != 150 {
				t.Errorf("Section %d White[%d] incorrect, got %d", section, i, val)
			}
		}
	}

	// Test invalid value
	err = bar.SetAllWhite(256)
	if err == nil {
		t.Error("Expected error for invalid value")
	}
}

func TestGetRGBWErrors(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	tests := []struct {
		name    string
		section int
		index   int
	}{
		{"Invalid section 0", 0, 0},
		{"Invalid section 3", 3, 0},
		{"Invalid index -1", 1, -1},
		{"Invalid index 6", 1, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, _, _, err := bar.GetRGBW(tt.section, tt.index)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}

func TestGetWhiteErrors(t *testing.T) {
	mock := mqtt.NewMockPublisher()
	bar, _ := NewLEDBar(0, mock, "test/topic")

	tests := []struct {
		name    string
		section int
		index   int
	}{
		{"Invalid section 0", 0, 0},
		{"Invalid section 3", 3, 0},
		{"Invalid index -1", 1, -1},
		{"Invalid index 13", 1, 13},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := bar.GetWhite(tt.section, tt.index)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}
