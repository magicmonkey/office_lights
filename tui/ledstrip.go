package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevin/office_lights/drivers/ledstrip"
)

// ledStripModel represents the LED strip section
type ledStripModel struct {
	driver        *ledstrip.LEDStrip
	activeControl int // 0=R, 1=G, 2=B
	r, g, b       int
}

func newLEDStripModel(driver *ledstrip.LEDStrip) ledStripModel {
	// Load current state from driver
	r, g, b := driver.GetColor()
	return ledStripModel{
		driver:        driver,
		activeControl: 0,
		r:             r,
		g:             g,
		b:             b,
	}
}

// Helper methods
func (m *ledStripModel) nextControl() {
	m.activeControl = (m.activeControl + 1) % 3
}

func (m *ledStripModel) prevControl() {
	m.activeControl = (m.activeControl - 1 + 3) % 3
}

func (m *ledStripModel) adjustValue(delta int) tea.Cmd {
	switch m.activeControl {
	case 0: // Red
		m.r = clamp(m.r+delta, 0, 255)
	case 1: // Green
		m.g = clamp(m.g+delta, 0, 255)
	case 2: // Blue
		m.b = clamp(m.b+delta, 0, 255)
	}
	return m.publish()
}

func (m *ledStripModel) publish() tea.Cmd {
	return func() tea.Msg {
		if err := m.driver.SetColor(m.r, m.g, m.b); err != nil {
			return publishErrorMsg{err}
		}
		return publishSuccessMsg{}
	}
}

// View renders this component
func (m ledStripModel) View(isActive bool) string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("LED Strip"))
	sb.WriteString("\n\n")

	// Red control
	if m.activeControl == 0 && isActive {
		sb.WriteString(activeControlStyle.Render("► R: "))
	} else {
		sb.WriteString(inactiveControlStyle.Render("  R: "))
	}
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.r)))
	sb.WriteString("\n")

	// Green control
	if m.activeControl == 1 && isActive {
		sb.WriteString(activeControlStyle.Render("► G: "))
	} else {
		sb.WriteString(inactiveControlStyle.Render("  G: "))
	}
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.g)))
	sb.WriteString("\n")

	// Blue control
	if m.activeControl == 2 && isActive {
		sb.WriteString(activeControlStyle.Render("► B: "))
	} else {
		sb.WriteString(inactiveControlStyle.Render("  B: "))
	}
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.b)))

	return sb.String()
}
