package web

import (
	"fmt"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
)

// RGBW represents a single RGBW LED
type RGBW struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	W int `json:"w"`
}

// LEDBarSection represents one section of the LED bar
type LEDBarSection struct {
	RGBW  []RGBW `json:"rgbw"`
	White []int  `json:"white"`
}

// LEDBarState represents the complete LED bar state
type LEDBarState struct {
	Section1 LEDBarSection `json:"section1"`
	Section2 LEDBarSection `json:"section2"`
}

// LEDStripState represents the LED strip state
type LEDStripState struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

// VideoLightState represents a video light state
type VideoLightState struct {
	On         bool `json:"on"`
	Brightness int  `json:"brightness"`
}

// State represents the complete system state
type State struct {
	LEDStrip    LEDStripState   `json:"ledStrip"`
	LEDBar      LEDBarState     `json:"ledBar"`
	VideoLight1 VideoLightState `json:"videoLight1"`
	VideoLight2 VideoLightState `json:"videoLight2"`
}

// BuildState reads current state from all drivers
func BuildState(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) (*State, error) {
	state := &State{}

	// LED Strip
	r, g, b := strip.GetColor()
	state.LEDStrip = LEDStripState{R: r, G: g, B: b}

	// LED Bar - initialize slices
	state.LEDBar.Section1.RGBW = make([]RGBW, 6)
	state.LEDBar.Section1.White = make([]int, 13)
	state.LEDBar.Section2.RGBW = make([]RGBW, 6)
	state.LEDBar.Section2.White = make([]int, 13)

	// Section 1 RGBW
	for i := 0; i < 6; i++ {
		r, g, b, w, err := bar.GetRGBW(1, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section1.RGBW[i] = RGBW{R: r, G: g, B: b, W: w}
	}

	// Section 1 White
	for i := 0; i < 13; i++ {
		val, err := bar.GetWhite(1, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section1.White[i] = val
	}

	// Section 2 RGBW
	for i := 0; i < 6; i++ {
		r, g, b, w, err := bar.GetRGBW(2, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section2.RGBW[i] = RGBW{R: r, G: g, B: b, W: w}
	}

	// Section 2 White
	for i := 0; i < 13; i++ {
		val, err := bar.GetWhite(2, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section2.White[i] = val
	}

	// Video Lights
	on1, brightness1 := vl1.GetState()
	state.VideoLight1 = VideoLightState{On: on1, Brightness: brightness1}

	on2, brightness2 := vl2.GetState()
	state.VideoLight2 = VideoLightState{On: on2, Brightness: brightness2}

	return state, nil
}

// ApplyState applies state to all drivers
func ApplyState(
	state *State,
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) error {
	// LED Strip
	if err := strip.SetColor(state.LEDStrip.R, state.LEDStrip.G, state.LEDStrip.B); err != nil {
		return fmt.Errorf("LED strip: %w", err)
	}

	// LED Bar Section 1 RGBW (use NoPublish to avoid multiple MQTT messages)
	for i, rgbw := range state.LEDBar.Section1.RGBW {
		if err := bar.SetRGBWNoPublish(1, i, rgbw.R, rgbw.G, rgbw.B, rgbw.W); err != nil {
			return fmt.Errorf("LED bar section1 RGBW[%d]: %w", i, err)
		}
	}

	// LED Bar Section 1 White
	for i, val := range state.LEDBar.Section1.White {
		if err := bar.SetWhiteNoPublish(1, i, val); err != nil {
			return fmt.Errorf("LED bar section1 white[%d]: %w", i, err)
		}
	}

	// LED Bar Section 2 RGBW
	for i, rgbw := range state.LEDBar.Section2.RGBW {
		if err := bar.SetRGBWNoPublish(2, i, rgbw.R, rgbw.G, rgbw.B, rgbw.W); err != nil {
			return fmt.Errorf("LED bar section2 RGBW[%d]: %w", i, err)
		}
	}

	// LED Bar Section 2 White
	for i, val := range state.LEDBar.Section2.White {
		if err := bar.SetWhiteNoPublish(2, i, val); err != nil {
			return fmt.Errorf("LED bar section2 white[%d]: %w", i, err)
		}
	}

	// Publish LED bar state once after all changes
	if err := bar.Publish(); err != nil {
		return fmt.Errorf("LED bar publish: %w", err)
	}

	// Video Light 1
	if state.VideoLight1.On {
		if err := vl1.TurnOn(state.VideoLight1.Brightness); err != nil {
			return fmt.Errorf("video light 1: %w", err)
		}
	} else {
		if err := vl1.TurnOff(); err != nil {
			return fmt.Errorf("video light 1: %w", err)
		}
	}

	// Video Light 2
	if state.VideoLight2.On {
		if err := vl2.TurnOn(state.VideoLight2.Brightness); err != nil {
			return fmt.Errorf("video light 2: %w", err)
		}
	} else {
		if err := vl2.TurnOff(); err != nil {
			return fmt.Errorf("video light 2: %w", err)
		}
	}

	return nil
}

// Validate checks if state values are within valid ranges
func (s *State) Validate() error {
	// LED Strip
	if s.LEDStrip.R < 0 || s.LEDStrip.R > 255 {
		return fmt.Errorf("LED strip R value out of range: %d", s.LEDStrip.R)
	}
	if s.LEDStrip.G < 0 || s.LEDStrip.G > 255 {
		return fmt.Errorf("LED strip G value out of range: %d", s.LEDStrip.G)
	}
	if s.LEDStrip.B < 0 || s.LEDStrip.B > 255 {
		return fmt.Errorf("LED strip B value out of range: %d", s.LEDStrip.B)
	}

	// LED Bar - validate RGBW values
	validateRGBW := func(section string, rgbwList []RGBW) error {
		if len(rgbwList) != 6 {
			return fmt.Errorf("LED bar %s RGBW must have 6 elements, got %d", section, len(rgbwList))
		}
		for i, rgbw := range rgbwList {
			if rgbw.R < 0 || rgbw.R > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] R out of range: %d", section, i, rgbw.R)
			}
			if rgbw.G < 0 || rgbw.G > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] G out of range: %d", section, i, rgbw.G)
			}
			if rgbw.B < 0 || rgbw.B > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] B out of range: %d", section, i, rgbw.B)
			}
			if rgbw.W < 0 || rgbw.W > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] W out of range: %d", section, i, rgbw.W)
			}
		}
		return nil
	}

	if err := validateRGBW("section1", s.LEDBar.Section1.RGBW); err != nil {
		return err
	}
	if err := validateRGBW("section2", s.LEDBar.Section2.RGBW); err != nil {
		return err
	}

	// LED Bar - validate white values
	validateWhite := func(section string, white []int) error {
		if len(white) != 13 {
			return fmt.Errorf("LED bar %s white must have 13 elements, got %d", section, len(white))
		}
		for i, val := range white {
			if val < 0 || val > 255 {
				return fmt.Errorf("LED bar %s white[%d] out of range: %d", section, i, val)
			}
		}
		return nil
	}

	if err := validateWhite("section1", s.LEDBar.Section1.White); err != nil {
		return err
	}
	if err := validateWhite("section2", s.LEDBar.Section2.White); err != nil {
		return err
	}

	// Video Lights
	if s.VideoLight1.Brightness < 0 || s.VideoLight1.Brightness > 100 {
		return fmt.Errorf("video light 1 brightness out of range: %d", s.VideoLight1.Brightness)
	}
	if s.VideoLight2.Brightness < 0 || s.VideoLight2.Brightness > 100 {
		return fmt.Errorf("video light 2 brightness out of range: %d", s.VideoLight2.Brightness)
	}

	return nil
}
