package streamdeck

import (
	"image"
	"sync"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
	sdlib "rafaelmartins.com/p/streamdeck"
)

// Mode represents the current operational mode of the Stream Deck interface
type Mode int

const (
	ModeLEDStrip Mode = iota
	ModeLEDBarRGBW
	ModeLEDBarWhite
	ModeVideoLights
)

// String returns the string representation of a Mode
func (m Mode) String() string {
	switch m {
	case ModeLEDStrip:
		return "LED Strip"
	case ModeLEDBarRGBW:
		return "LED Bar RGBW"
	case ModeLEDBarWhite:
		return "LED Bar White"
	case ModeVideoLights:
		return "Video Lights"
	default:
		return "Unknown"
	}
}

// SectionData represents data to display in one touchscreen section
type SectionData struct {
	Label string // e.g., "Red", "Green", "Light1"
	Value int    // 0-255 or 0-100
	Active bool  // Whether this section is active in the current mode
}

// StreamDeckUI manages the Stream Deck+ interface
type StreamDeckUI struct {
	device      *sdlib.Device
	ledStrip    *ledstrip.LEDStrip
	ledBar      *ledbar.LEDBar
	videoLight1 *videolight.VideoLight
	videoLight2 *videolight.VideoLight

	mu          sync.Mutex
	currentMode Mode
	lastValues  [4]int // Store last non-zero values for toggle functionality

	// Cached images
	buttonImages [8]image.Image
	touchImage   image.Image

	// Control
	quit chan struct{}
}

// NewStreamDeckUI creates a new Stream Deck UI instance
func NewStreamDeckUI(
	ledStrip *ledstrip.LEDStrip,
	ledBar *ledbar.LEDBar,
	videoLight1 *videolight.VideoLight,
	videoLight2 *videolight.VideoLight,
) (*StreamDeckUI, error) {
	// Find Stream Deck devices
	devices, err := sdlib.Enumerate()
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, ErrNoDevice
	}

	// Use the first device found
	device := devices[0]

	// Open the device
	if err := device.Open(); err != nil {
		return nil, err
	}

	ui := &StreamDeckUI{
		device:      device,
		ledStrip:    ledStrip,
		ledBar:      ledBar,
		videoLight1: videoLight1,
		videoLight2: videoLight2,
		currentMode: ModeLEDStrip, // Default mode
		quit:        make(chan struct{}),
	}

	return ui, nil
}

// Close cleans up the Stream Deck UI
func (s *StreamDeckUI) Close() error {
	close(s.quit)
	if s.device != nil {
		return s.device.Close()
	}
	return nil
}
