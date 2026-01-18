package streamdeck

import (
	"log"
	"time"

	sdlib "rafaelmartins.com/p/streamdeck"
)

// Run starts the Stream Deck UI event loop
func (s *StreamDeckUI) Run() error {
	// Initialize the display
	if err := s.initializeDisplay(); err != nil {
		return err
	}

	// Register event handlers
	if err := s.registerHandlers(); err != nil {
		return err
	}

	// Start periodic touchscreen updates
	ticker := time.NewTicker(100 * time.Millisecond) // 10 FPS
	defer ticker.Stop()

	// Start listening for events in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := s.device.Listen(errCh); err != nil {
			log.Printf("Stream Deck Listen error: %v", err)
		}
	}()

	log.Println("Stream Deck UI running")

	for {
		select {
		case <-s.quit:
			log.Println("Stream Deck UI shutting down")
			return nil

		case err := <-errCh:
			if err != nil {
				log.Printf("Stream Deck error: %v", err)
			}

		case <-ticker.C:
			// Update touchscreen display
			if err := s.updateTouchscreen(); err != nil {
				log.Printf("Error updating touchscreen: %v", err)
			}
		}
	}
}

// initializeDisplay sets up the initial button images and touchscreen
func (s *StreamDeckUI) initializeDisplay() error {
	log.Println("Initializing Stream Deck display...")

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

// registerHandlers sets up event handlers for keys, dials, and touchscreen
func (s *StreamDeckUI) registerHandlers() error {
	// Register key handlers (buttons 1-8)
	keyIDs := []sdlib.KeyID{
		sdlib.KEY_1, sdlib.KEY_2, sdlib.KEY_3, sdlib.KEY_4,
		sdlib.KEY_5, sdlib.KEY_6, sdlib.KEY_7, sdlib.KEY_8,
	}

	for i, keyID := range keyIDs {
		index := i // Capture loop variable
		if err := s.device.AddKeyHandler(keyID, func(d *sdlib.Device, k *sdlib.Key) error {
			s.handleButtonPress(index)
			return nil
		}); err != nil {
			return err
		}
	}

	// Register dial rotate handlers (dials 1-4)
	dialIDs := []sdlib.DialID{sdlib.DIAL_1, sdlib.DIAL_2, sdlib.DIAL_3, sdlib.DIAL_4}
	for i, dialID := range dialIDs {
		index := i // Capture loop variable
		if err := s.device.AddDialRotateHandler(dialID, func(d *sdlib.Device, di *sdlib.Dial, delta int8) error {
			s.handleDialRotate(index, int(delta))
			return nil
		}); err != nil {
			// Dial might not be supported on this device, log but don't fail
			log.Printf("Warning: Could not register dial %d handler: %v", i+1, err)
		}
	}

	// Register dial switch handlers (dial clicks)
	for i, dialID := range dialIDs {
		index := i // Capture loop variable
		if err := s.device.AddDialSwitchHandler(dialID, func(d *sdlib.Device, di *sdlib.Dial) error {
			s.handleDialPress(index)
			return nil
		}); err != nil {
			// Dial might not be supported on this device, log but don't fail
			log.Printf("Warning: Could not register dial switch %d handler: %v", i+1, err)
		}
	}

	return nil
}
