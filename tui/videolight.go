package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevin/office_lights/drivers/videolight"
)

// videoLightModel represents a video light section
type videoLightModel struct {
	driver        *videolight.VideoLight
	lightID       int // 1 or 2
	activeControl int // 0=on/off, 1=brightness
	on            bool
	brightness    int
}

func newVideoLightModel(driver *videolight.VideoLight, lightID int) videoLightModel {
	// Load current state from driver
	on, brightness := driver.GetState()
	return videoLightModel{
		driver:        driver,
		lightID:       lightID,
		activeControl: 0,
		on:            on,
		brightness:    brightness,
	}
}

func (m *videoLightModel) nextControl() {
	m.activeControl = (m.activeControl + 1) % 2
}

func (m *videoLightModel) prevControl() {
	m.activeControl = (m.activeControl - 1 + 2) % 2
}

// refresh updates the model's values from the driver
func (m *videoLightModel) refresh() {
	m.on, m.brightness = m.driver.GetState()
}

func (m *videoLightModel) adjustValue(delta int) tea.Cmd {
	if m.activeControl == 0 {
		// On/off toggle - ignore delta for this control
		return nil
	}
	// Brightness
	m.brightness = clamp(m.brightness+delta, 0, 100)
	return m.publish()
}

func (m *videoLightModel) toggle() tea.Cmd {
	if m.activeControl == 0 {
		m.on = !m.on
	}
	return m.publish()
}

func (m *videoLightModel) publish() tea.Cmd {
	return func() tea.Msg {
		var err error
		if m.on {
			err = m.driver.TurnOn(m.brightness)
		} else {
			err = m.driver.TurnOff()
		}
		if err != nil {
			return publishErrorMsg{err}
		}
		return publishSuccessMsg{}
	}
}

// View renders this component
func (m videoLightModel) View(isActive bool) string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render(fmt.Sprintf("Video Light %d", m.lightID)))
	sb.WriteString("\n\n")

	// On/Off control
	if m.activeControl == 0 && isActive {
		sb.WriteString(activeControlStyle.Render("► On: "))
	} else {
		sb.WriteString(inactiveControlStyle.Render("  On: "))
	}
	onStr := "false"
	if m.on {
		onStr = "true"
	}
	sb.WriteString(valueStyle.Render(onStr))
	sb.WriteString("\n")

	// Brightness control
	if m.activeControl == 1 && isActive {
		sb.WriteString(activeControlStyle.Render("► Brightness: "))
	} else {
		sb.WriteString(inactiveControlStyle.Render("  Brightness: "))
	}
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.brightness)))

	return sb.String()
}
