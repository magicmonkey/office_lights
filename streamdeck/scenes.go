package streamdeck

import (
	"log"

	"github.com/kevin/office_lights/storage"
)

// saveScene captures current light state and saves to database
func (s *StreamDeckUI) saveScene(slotIndex int) {
	log.Printf("Saving scene %d...", slotIndex+1)

	// Gather current state from all drivers
	data := &storage.SceneData{
		LEDStrip: storage.LEDStripState{
			Red:   s.ledStrip.R(),
			Green: s.ledStrip.G(),
			Blue:  s.ledStrip.B(),
		},
	}

	// Gather LED bar state
	channels := s.ledBar.GetChannels()
	for i, value := range channels {
		data.LEDBarLEDs = append(data.LEDBarLEDs, storage.LEDBarLEDState{
			LEDBarID:   s.ledBar.GetBarID(),
			ChannelNum: i,
			Value:      value,
		})
	}

	// Gather video light states
	data.VideoLights = []storage.VideoLightState{
		{
			ID:         0,
			On:         s.videoLight1.IsOn(),
			Brightness: s.videoLight1.Brightness(),
		},
		{
			ID:         1,
			On:         s.videoLight2.IsOn(),
			Brightness: s.videoLight2.Brightness(),
		},
	}

	// Save to database
	if err := s.storage.SaveScene(slotIndex, data); err != nil {
		log.Printf("Error saving scene %d: %v", slotIndex+1, err)
		return
	}

	log.Printf("Scene %d saved successfully", slotIndex+1)

	// Update display to show saved status
	if err := s.updateButtons(); err != nil {
		log.Printf("Error updating buttons: %v", err)
	}
	if err := s.updateTouchscreen(); err != nil {
		log.Printf("Error updating touchscreen: %v", err)
	}
}

// recallScene loads saved state and applies to all lights
func (s *StreamDeckUI) recallScene(slotIndex int) {
	log.Printf("Recalling scene %d...", slotIndex+1)

	// Load from database
	data, err := s.storage.LoadScene(slotIndex)
	if err != nil {
		log.Printf("Error loading scene %d: %v", slotIndex+1, err)
		return
	}

	if data == nil {
		log.Printf("Scene %d is empty", slotIndex+1)
		return
	}

	// Apply to LED strip
	if err := s.ledStrip.SetColor(data.LEDStrip.Red, data.LEDStrip.Green, data.LEDStrip.Blue); err != nil {
		log.Printf("Error setting LED strip: %v", err)
	}

	// Apply to LED bar
	if len(data.LEDBarLEDs) > 0 {
		channels := make([]int, 77)
		for _, led := range data.LEDBarLEDs {
			if led.ChannelNum >= 0 && led.ChannelNum < 77 {
				channels[led.ChannelNum] = led.Value
			}
		}
		if err := s.ledBar.SetChannels(channels); err != nil {
			log.Printf("Error setting LED bar: %v", err)
		}
	}

	// Apply to video lights
	for _, vl := range data.VideoLights {
		switch vl.ID {
		case 0:
			if vl.On {
				if err := s.videoLight1.TurnOn(vl.Brightness); err != nil {
					log.Printf("Error setting video light 1: %v", err)
				}
			} else {
				if err := s.videoLight1.TurnOff(); err != nil {
					log.Printf("Error turning off video light 1: %v", err)
				}
			}
		case 1:
			if vl.On {
				if err := s.videoLight2.TurnOn(vl.Brightness); err != nil {
					log.Printf("Error setting video light 2: %v", err)
				}
			} else {
				if err := s.videoLight2.TurnOff(); err != nil {
					log.Printf("Error turning off video light 2: %v", err)
				}
			}
		}
	}

	log.Printf("Scene %d recalled successfully", slotIndex+1)
}
