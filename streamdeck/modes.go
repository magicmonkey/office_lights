package streamdeck

// getSectionData returns the section data for the current mode
func (s *StreamDeckUI) getSectionData() [4]SectionData {
	var sections [4]SectionData

	switch s.currentMode {
	case ModeLEDStrip:
		sections = s.getLEDStripSections()
	case ModeLEDBarRGBW:
		sections = s.getLEDBarRGBWSections()
	case ModeLEDBarWhite:
		sections = s.getLEDBarWhiteSections()
	case ModeVideoLights:
		sections = s.getVideoLightsSections()
	}

	return sections
}

// getLEDStripSections returns section data for LED Strip mode (3 active sections)
func (s *StreamDeckUI) getLEDStripSections() [4]SectionData {
	r, g, b := s.ledStrip.R(), s.ledStrip.G(), s.ledStrip.B()

	return [4]SectionData{
		{Label: "Red", Value: r, Active: true},
		{Label: "Green", Value: g, Active: true},
		{Label: "Blue", Value: b, Active: true},
		{Label: "", Value: 0, Active: false},
	}
}

// getLEDBarRGBWSections returns section data for LED Bar RGBW mode (4 active sections)
// Shows values for the first RGBW LED in section 1
func (s *StreamDeckUI) getLEDBarRGBWSections() [4]SectionData {
	// Get RGBW values from the first LED (index 0) in section 1
	r, g, b, w, err := s.ledBar.GetRGBW(1, 0)
	if err != nil {
		// Return zeros on error
		r, g, b, w = 0, 0, 0, 0
	}

	return [4]SectionData{
		{Label: "Red", Value: r, Active: true},
		{Label: "Green", Value: g, Active: true},
		{Label: "Blue", Value: b, Active: true},
		{Label: "White", Value: w, Active: true},
	}
}

// getLEDBarWhiteSections returns section data for LED Bar White mode (2 active sections)
// Shows average brightness for all white LEDs in each section
func (s *StreamDeckUI) getLEDBarWhiteSections() [4]SectionData {
	// Calculate average brightness for each section's white LEDs
	section1Avg := s.ledBar.GetAverageWhite(1)
	section2Avg := s.ledBar.GetAverageWhite(2)

	return [4]SectionData{
		{Label: "Section 1", Value: section1Avg, Active: true},
		{Label: "Section 2", Value: section2Avg, Active: true},
		{Label: "", Value: 0, Active: false},
		{Label: "", Value: 0, Active: false},
	}
}

// getVideoLightsSections returns section data for Video Lights mode (4 active sections)
// Dials 0 and 1 are coarse adjustments (±5), dials 2 and 3 are fine-tune (±1)
func (s *StreamDeckUI) getVideoLightsSections() [4]SectionData {
	on1, brightness1 := s.videoLight1.IsOn(), s.videoLight1.Brightness()
	on2, brightness2 := s.videoLight2.IsOn(), s.videoLight2.Brightness()

	label1 := "Light 1"
	if !on1 {
		label1 = "Light 1 (OFF)"
	}

	label2 := "Light 2"
	if !on2 {
		label2 = "Light 2 (OFF)"
	}

	fineTuneLabel1 := "L1 Fine"
	if !on1 {
		fineTuneLabel1 = "L1 Fine (OFF)"
	}

	fineTuneLabel2 := "L2 Fine"
	if !on2 {
		fineTuneLabel2 = "L2 Fine (OFF)"
	}

	return [4]SectionData{
		{Label: label1, Value: brightness1, Active: true},
		{Label: label2, Value: brightness2, Active: true},
		{Label: fineTuneLabel1, Value: brightness1, Active: true},
		{Label: fineTuneLabel2, Value: brightness2, Active: true},
	}
}
