package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
)

// Section represents which light section is active
type Section int

const (
	SectionLEDStrip Section = iota
	SectionLEDBar
	SectionVideoLight1
	SectionVideoLight2
	SectionCount
)

// Model is the root Bubbletea model
type Model struct {
	// Focus management
	activeSection Section

	// Component models
	ledStrip    ledStripModel
	ledBar      ledBarModel
	videoLight1 videoLightModel
	videoLight2 videoLightModel

	// Driver references
	stripDriver *ledstrip.LEDStrip
	barDriver   *ledbar.LEDBar
	vl1Driver   *videolight.VideoLight
	vl2Driver   *videolight.VideoLight

	// UI state
	width  int
	height int
	ready  bool
	err    error
}

// New creates a new TUI model
func New(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) Model {
	return Model{
		activeSection: SectionLEDStrip,
		stripDriver:   strip,
		barDriver:     bar,
		vl1Driver:     vl1,
		vl2Driver:     vl2,
		ledStrip:      newLEDStripModel(strip),
		ledBar:        newLEDBarModel(bar),
		videoLight1:   newVideoLightModel(vl1, 1),
		videoLight2:   newVideoLightModel(vl2, 2),
	}
}

// Init initializes the model (Bubbletea requirement)
func (m Model) Init() tea.Cmd {
	// Start the periodic refresh tick
	return tick()
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
