package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
)

// Run starts the TUI
func Run(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) error {
	m := New(strip, bar, vl1, vl2)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI error: %w", err)
	}

	return nil
}
