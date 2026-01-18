package streamdeck

import "log"

const (
	dialIncrement = 5 // Amount to increment/decrement per dial tick
)

// handleButtonPress processes button press events
func (s *StreamDeckUI) handleButtonPress(buttonIndex int) {
	log.Printf("Button %d pressed", buttonIndex)

	// Buttons 0-3 (top row): Mode selection
	if buttonIndex < 4 {
		newMode := Mode(buttonIndex)
		if newMode != s.currentMode {
			log.Printf("Switching mode from %s to %s", s.currentMode, newMode)
			s.currentMode = newMode
			// Update button display to reflect new mode
			if err := s.updateButtons(); err != nil {
				log.Printf("Error updating buttons: %v", err)
			}
			// Update touchscreen immediately to show new mode data
			if err := s.updateTouchscreen(); err != nil {
				log.Printf("Error updating touchscreen: %v", err)
			}
		}
	}

	// Buttons 4-7 (bottom row): Reserved/unused
}

// handleDialRotate processes dial rotation events
func (s *StreamDeckUI) handleDialRotate(dialIndex int, ticks int) {
	if dialIndex < 0 || dialIndex > 3 {
		log.Printf("Invalid dial index: %d", dialIndex)
		return
	}

	log.Printf("Dial %d rotated: %d ticks", dialIndex, ticks)

	// Get section data to determine if this dial is active
	sections := s.getSectionData()
	if !sections[dialIndex].Active {
		log.Printf("Dial %d is inactive in current mode", dialIndex)
		return
	}

	// Calculate increment
	increment := ticks * dialIncrement

	// Apply increment based on current mode and dial index
	switch s.currentMode {
	case ModeLEDStrip:
		s.adjustLEDStrip(dialIndex, increment)
	case ModeLEDBarRGBW:
		s.adjustLEDBarRGBW(dialIndex, increment)
	case ModeLEDBarWhite:
		s.adjustLEDBarWhite(dialIndex, increment)
	case ModeVideoLights:
		s.adjustVideoLights(dialIndex, increment)
	}
}

// handleDialPress processes dial click/press events
func (s *StreamDeckUI) handleDialPress(dialIndex int) {
	if dialIndex < 0 || dialIndex > 3 {
		log.Printf("Invalid dial index: %d", dialIndex)
		return
	}

	log.Printf("Dial %d pressed", dialIndex)

	// Get section data to determine if this dial is active
	sections := s.getSectionData()
	if !sections[dialIndex].Active {
		log.Printf("Dial %d is inactive in current mode", dialIndex)
		return
	}

	// Toggle based on current mode and dial index
	switch s.currentMode {
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

// handleTouch processes touchscreen touch events
func (s *StreamDeckUI) handleTouch(x, y int) {
	// Calculate which section was touched
	section := x / sectionWidth
	log.Printf("Touchscreen touched at (%d, %d) -> section %d", x, y, section)

	// Optional: provide visual feedback or additional functionality
	// For now, we just log the event
}

// LED Strip adjustment functions

func (s *StreamDeckUI) adjustLEDStrip(dialIndex int, increment int) {
	r, g, b := s.ledStrip.R(), s.ledStrip.G(), s.ledStrip.B()

	switch dialIndex {
	case 0: // Red
		r = clamp(r + increment)
		s.lastValues[0] = r
	case 1: // Green
		g = clamp(g + increment)
		s.lastValues[1] = g
	case 2: // Blue
		b = clamp(b + increment)
		s.lastValues[2] = b
	}

	if err := s.ledStrip.SetColor(r, g, b); err != nil {
		log.Printf("Error setting LED strip color: %v", err)
	}
}

func (s *StreamDeckUI) toggleLEDStrip(dialIndex int) {
	r, g, b := s.ledStrip.R(), s.ledStrip.G(), s.ledStrip.B()

	switch dialIndex {
	case 0: // Red
		if r == 0 {
			r = s.lastValues[0]
			if r == 0 {
				r = 255
			}
		} else {
			s.lastValues[0] = r
			r = 0
		}
	case 1: // Green
		if g == 0 {
			g = s.lastValues[1]
			if g == 0 {
				g = 255
			}
		} else {
			s.lastValues[1] = g
			g = 0
		}
	case 2: // Blue
		if b == 0 {
			b = s.lastValues[2]
			if b == 0 {
				b = 255
			}
		} else {
			s.lastValues[2] = b
			b = 0
		}
	}

	if err := s.ledStrip.SetColor(r, g, b); err != nil {
		log.Printf("Error setting LED strip color: %v", err)
	}
}

// LED Bar RGBW adjustment functions

func (s *StreamDeckUI) adjustLEDBarRGBW(dialIndex int, increment int) {
	// Get current values from first LED in section 1
	r, g, b, w, err := s.ledBar.GetRGBW(1, 0)
	if err != nil {
		log.Printf("Error getting LED bar RGBW: %v", err)
		return
	}

	switch dialIndex {
	case 0: // Red
		r = clamp(r + increment)
	case 1: // Green
		g = clamp(g + increment)
	case 2: // Blue
		b = clamp(b + increment)
	case 3: // White
		w = clamp(w + increment)
	}

	// Set all RGBW LEDs in both sections to the same value
	if err := s.ledBar.SetAllRGBW(r, g, b, w); err != nil {
		log.Printf("Error setting LED bar RGBW: %v", err)
	}
}

func (s *StreamDeckUI) toggleLEDBarRGBW(dialIndex int) {
	r, g, b, w, err := s.ledBar.GetRGBW(1, 0)
	if err != nil {
		log.Printf("Error getting LED bar RGBW: %v", err)
		return
	}

	switch dialIndex {
	case 0: // Red
		if r == 0 {
			r = s.lastValues[0]
			if r == 0 {
				r = 255
			}
		} else {
			s.lastValues[0] = r
			r = 0
		}
	case 1: // Green
		if g == 0 {
			g = s.lastValues[1]
			if g == 0 {
				g = 255
			}
		} else {
			s.lastValues[1] = g
			g = 0
		}
	case 2: // Blue
		if b == 0 {
			b = s.lastValues[2]
			if b == 0 {
				b = 255
			}
		} else {
			s.lastValues[2] = b
			b = 0
		}
	case 3: // White
		if w == 0 {
			w = s.lastValues[3]
			if w == 0 {
				w = 255
			}
		} else {
			s.lastValues[3] = w
			w = 0
		}
	}

	if err := s.ledBar.SetAllRGBW(r, g, b, w); err != nil {
		log.Printf("Error setting LED bar RGBW: %v", err)
	}
}

// LED Bar White adjustment functions

func (s *StreamDeckUI) adjustLEDBarWhite(dialIndex int, increment int) {
	section1Avg := s.ledBar.GetAverageWhite(1)
	section2Avg := s.ledBar.GetAverageWhite(2)

	switch dialIndex {
	case 0: // Section 1
		newValue := clamp(section1Avg + increment)
		if err := s.ledBar.SetAllWhite(1, newValue); err != nil {
			log.Printf("Error setting LED bar section 1 white: %v", err)
		}
	case 1: // Section 2
		newValue := clamp(section2Avg + increment)
		if err := s.ledBar.SetAllWhite(2, newValue); err != nil {
			log.Printf("Error setting LED bar section 2 white: %v", err)
		}
	}
}

func (s *StreamDeckUI) toggleLEDBarWhite(dialIndex int) {
	section1Avg := s.ledBar.GetAverageWhite(1)
	section2Avg := s.ledBar.GetAverageWhite(2)

	switch dialIndex {
	case 0: // Section 1
		newValue := 0
		if section1Avg == 0 {
			newValue = s.lastValues[0]
			if newValue == 0 {
				newValue = 255
			}
		} else {
			s.lastValues[0] = section1Avg
		}
		if err := s.ledBar.SetAllWhite(1, newValue); err != nil {
			log.Printf("Error setting LED bar section 1 white: %v", err)
		}
	case 1: // Section 2
		newValue := 0
		if section2Avg == 0 {
			newValue = s.lastValues[1]
			if newValue == 0 {
				newValue = 255
			}
		} else {
			s.lastValues[1] = section2Avg
		}
		if err := s.ledBar.SetAllWhite(2, newValue); err != nil {
			log.Printf("Error setting LED bar section 2 white: %v", err)
		}
	}
}

// Video Lights adjustment functions

func (s *StreamDeckUI) adjustVideoLights(dialIndex int, increment int) {
	switch dialIndex {
	case 0: // Video Light 1 (coarse adjustment)
		s.adjustVideoLight1(increment)
	case 1: // Video Light 2 (coarse adjustment)
		s.adjustVideoLight2(increment)
	case 2: // Video Light 1 (fine-tune, increment of 1 per tick)
		fineIncrement := increment / dialIncrement // Convert back to ticks for ±1 adjustment
		s.adjustVideoLight1(fineIncrement)
	case 3: // Video Light 2 (fine-tune, increment of 1 per tick)
		fineIncrement := increment / dialIncrement // Convert back to ticks for ±1 adjustment
		s.adjustVideoLight2(fineIncrement)
	}
}

func (s *StreamDeckUI) adjustVideoLight1(increment int) {
	on := s.videoLight1.IsOn()
	if !on && increment > 0 {
		// Light is off and dial turning up: turn on at brightness 0, then nudge up
		brightness := clamp100(increment)
		if err := s.videoLight1.TurnOn(brightness); err != nil {
			log.Printf("Error turning on video light 1: %v", err)
		}
	} else if on {
		brightness := clamp100(s.videoLight1.Brightness() + increment)
		if err := s.videoLight1.TurnOn(brightness); err != nil {
			log.Printf("Error setting video light 1 brightness: %v", err)
		}
	}
}

func (s *StreamDeckUI) adjustVideoLight2(increment int) {
	on := s.videoLight2.IsOn()
	if !on && increment > 0 {
		// Light is off and dial turning up: turn on at brightness 0, then nudge up
		brightness := clamp100(increment)
		if err := s.videoLight2.TurnOn(brightness); err != nil {
			log.Printf("Error turning on video light 2: %v", err)
		}
	} else if on {
		brightness := clamp100(s.videoLight2.Brightness() + increment)
		if err := s.videoLight2.TurnOn(brightness); err != nil {
			log.Printf("Error setting video light 2 brightness: %v", err)
		}
	}
}

func (s *StreamDeckUI) toggleVideoLights(dialIndex int) {
	switch dialIndex {
	case 0, 2: // Video Light 1 (dial 0 = coarse, dial 2 = fine-tune)
		if s.videoLight1.IsOn() {
			if err := s.videoLight1.TurnOff(); err != nil {
				log.Printf("Error turning off video light 1: %v", err)
			}
		} else {
			brightness := s.videoLight1.Brightness()
			if brightness == 0 {
				brightness = 100 // Default brightness
			}
			if err := s.videoLight1.TurnOn(brightness); err != nil {
				log.Printf("Error turning on video light 1: %v", err)
			}
		}
	case 1, 3: // Video Light 2 (dial 1 = coarse, dial 3 = fine-tune)
		if s.videoLight2.IsOn() {
			if err := s.videoLight2.TurnOff(); err != nil {
				log.Printf("Error turning off video light 2: %v", err)
			}
		} else {
			brightness := s.videoLight2.Brightness()
			if brightness == 0 {
				brightness = 100 // Default brightness
			}
			if err := s.videoLight2.TurnOn(brightness); err != nil {
				log.Printf("Error turning on video light 2: %v", err)
			}
		}
	}
}

// Helper function to clamp values to 0-255 range
func clamp(value int) int {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

// Helper function to clamp values to 0-100 range (for video lights)
func clamp100(value int) int {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}
