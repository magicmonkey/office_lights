package streamdeck

import (
	"image"
	"sync"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
	sdlib "rafaelmartins.com/p/streamdeck"
)

// Tab represents the currently selected tab on the Stream Deck
type Tab int

const (
	TabLightControl Tab = iota // Tab 1: Light control (existing functionality)
	TabFuture2                 // Tab 2: Reserved for future use
	TabFuture3                 // Tab 3: Reserved for future use
	TabFuture4                 // Tab 4: Reserved for future use
)

// String returns the string representation of a Tab
func (t Tab) String() string {
	switch t {
	case TabLightControl:
		return "Lights"
	case TabFuture2:
		return "Tab 2"
	case TabFuture3:
		return "Tab 3"
	case TabFuture4:
		return "Tab 4"
	default:
		return "Unknown"
	}
}

// Mode represents the current operational mode within Tab 1 (Light Control)
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
	Label    string // e.g., "Red", "Green", "Light1"
	Value    int    // Current value
	MaxValue int    // Maximum value (255 for LEDs, 100 for video lights)
	Active   bool   // Whether this section is active in the current mode
}

// StreamDeckUI manages the Stream Deck+ interface
type StreamDeckUI struct {
	device      *sdlib.Device
	ledStrip    *ledstrip.LEDStrip
	ledBar      *ledbar.LEDBar
	videoLight1 *videolight.VideoLight
	videoLight2 *videolight.VideoLight

	mu          sync.Mutex
	currentTab  Tab  // Currently selected tab (0-3)
	currentMode Mode // Mode within TabLightControl
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
		currentTab:  TabLightControl, // Default to Light Control tab
		currentMode: ModeLEDStrip,    // Default mode within Light Control
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
