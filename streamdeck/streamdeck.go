package streamdeck

import (
	"log"
	"time"

	sdlib "github.com/muesli/streamdeck"
)

// Run starts the Stream Deck UI event loop
func (s *StreamDeckUI) Run() error {
	// Initialize the display
	if err := s.initializeDisplay(); err != nil {
		return err
	}

	// Start periodic touchscreen updates
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
	defer ticker.Stop()

	// Listen for key events
	keyCh, err := s.device.ReadKeys()
	if err != nil {
		return err
	}

	log.Println("Stream Deck UI running")

	for {
		select {
		case <-s.quit:
			log.Println("Stream Deck UI shutting down")
			return nil

		case <-ticker.C:
			// Update touchscreen display (currently not functional - library limitation)
			if err := s.updateTouchscreen(); err != nil {
				log.Printf("Error updating touchscreen: %v", err)
			}

		case key, ok := <-keyCh:
			if !ok {
				log.Println("Stream Deck key channel closed")
				return nil
			}
			s.handleEvent(key)
		}
	}
}

// initializeDisplay sets up the initial button images and touchscreen
func (s *StreamDeckUI) initializeDisplay() error {
	log.Println("Initializing Stream Deck display...")

	// Clear all buttons
	if err := s.device.Clear(); err != nil {
		return err
	}

	// Render and set button images
	if err := s.updateButtons(); err != nil {
		return err
	}

	// Render and set touchscreen
	if err := s.updateTouchscreen(); err != nil {
		return err
	}

	return nil
}

// handleEvent processes Stream Deck events
// Note: The github.com/muesli/streamdeck library only supports key press events.
// Dial and touchscreen events are not available with this library.
func (s *StreamDeckUI) handleEvent(key sdlib.Key) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Only handle button presses
	if key.Pressed {
		s.handleButtonPress(int(key.Index))
	}
}
