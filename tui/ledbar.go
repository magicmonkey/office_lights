package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kevin/office_lights/drivers/ledbar"
)

// ledBarModel represents the LED bar section
type ledBarModel struct {
	driver        *ledbar.LEDBar
	mode          int // 0=RGBW, 1=White
	section       int // 1 or 2
	activeControl int // Depends on mode

	// RGBW mode controls: 0=mode, 1=section, 2=index, 3=R, 4=G, 5=B, 6=W
	rgbwIndex int // 0-5 (LED 1-6)
	r, g, b, w int

	// White mode controls: 0=mode, 1=section, 2=index, 3=brightness
	whiteIndex      int // 0-12 (LED 1-13)
	whiteBrightness int
}

func newLEDBarModel(driver *ledbar.LEDBar) ledBarModel {
	// Initialize with default values
	// Try to load current state from driver for first RGBW LED
	r, g, b, w, err := driver.GetRGBW(1, 0)
	if err != nil {
		r, g, b, w = 0, 0, 0, 0
	}

	brightness, err := driver.GetWhite(1, 0)
	if err != nil {
		brightness = 0
	}

	return ledBarModel{
		driver:          driver,
		mode:            0, // Start in RGBW mode
		section:         1, // Start with section 1
		activeControl:   0,
		rgbwIndex:       0,
		r:               r,
		g:               g,
		b:               b,
		w:               w,
		whiteIndex:      0,
		whiteBrightness: brightness,
	}
}

func (m *ledBarModel) nextControl() {
	if m.mode == 0 { // RGBW mode
		m.activeControl = (m.activeControl + 1) % 7
	} else { // White mode
		m.activeControl = (m.activeControl + 1) % 4
	}
}

func (m *ledBarModel) prevControl() {
	if m.mode == 0 {
		m.activeControl = (m.activeControl - 1 + 7) % 7
	} else {
		m.activeControl = (m.activeControl - 1 + 4) % 4
	}
}

// refresh updates the model's values from the driver
func (m *ledBarModel) refresh() {
	// Refresh RGBW values for current selection
	r, g, b, w, err := m.driver.GetRGBW(m.section, m.rgbwIndex)
	if err == nil {
		m.r, m.g, m.b, m.w = r, g, b, w
	}

	// Refresh white brightness for current selection
	brightness, err := m.driver.GetWhite(m.section, m.whiteIndex)
	if err == nil {
		m.whiteBrightness = brightness
	}
}

func (m *ledBarModel) adjustValue(delta int) tea.Cmd {
	if m.mode == 0 { // RGBW mode
		switch m.activeControl {
		case 0: // Mode - cycle between RGBW and White
			if delta != 0 {
				m.mode = 1
				m.activeControl = 0
			}
		case 1: // Section
			if delta > 0 {
				m.section = 2
			} else if delta < 0 {
				m.section = 1
			}
		case 2: // RGBW Index
			m.rgbwIndex = clamp(m.rgbwIndex+delta, 0, 5)
			// Load the values for this LED
			r, g, b, w, err := m.driver.GetRGBW(m.section, m.rgbwIndex)
			if err == nil {
				m.r, m.g, m.b, m.w = r, g, b, w
			}
		case 3: // R
			m.r = clamp(m.r+delta, 0, 255)
			return m.publish()
		case 4: // G
			m.g = clamp(m.g+delta, 0, 255)
			return m.publish()
		case 5: // B
			m.b = clamp(m.b+delta, 0, 255)
			return m.publish()
		case 6: // W
			m.w = clamp(m.w+delta, 0, 255)
			return m.publish()
		}
	} else { // White mode
		switch m.activeControl {
		case 0: // Mode - cycle between RGBW and White
			if delta != 0 {
				m.mode = 0
				m.activeControl = 0
			}
		case 1: // Section
			if delta > 0 {
				m.section = 2
			} else if delta < 0 {
				m.section = 1
			}
		case 2: // White Index
			m.whiteIndex = clamp(m.whiteIndex+delta, 0, 12)
			// Load the value for this LED
			brightness, err := m.driver.GetWhite(m.section, m.whiteIndex)
			if err == nil {
				m.whiteBrightness = brightness
			}
		case 3: // Brightness
			m.whiteBrightness = clamp(m.whiteBrightness+delta, 0, 255)
			return m.publish()
		}
	}
	return nil
}

func (m *ledBarModel) publish() tea.Cmd {
	return func() tea.Msg {
		var err error
		if m.mode == 0 {
			err = m.driver.SetRGBW(m.section, m.rgbwIndex, m.r, m.g, m.b, m.w)
		} else {
			err = m.driver.SetWhite(m.section, m.whiteIndex, m.whiteBrightness)
		}
		if err != nil {
			return publishErrorMsg{err}
		}
		return publishSuccessMsg{}
	}
}

// View renders this component
func (m ledBarModel) View(isActive bool) string {
	var sb strings.Builder

	sb.WriteString(titleStyle.Render("LED Bar"))
	sb.WriteString("\n\n")

	// Mode selector
	if m.activeControl == 0 && isActive {
		sb.WriteString(activeControlStyle.Render("► Mode: "))
	} else {
		sb.WriteString(inactiveControlStyle.Render("  Mode: "))
	}
	modeStr := "RGBW"
	if m.mode == 1 {
		modeStr = "White"
	}
	sb.WriteString(valueStyle.Render(modeStr))
	sb.WriteString("\n")

	// Section selector
	if m.activeControl == 1 && isActive {
		sb.WriteString(activeControlStyle.Render("► Section: "))
	} else {
		sb.WriteString(inactiveControlStyle.Render("  Section: "))
	}
	sb.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.section)))
	sb.WriteString("\n")

	if m.mode == 0 { // RGBW mode
		// LED Index
		if m.activeControl == 2 && isActive {
			sb.WriteString(activeControlStyle.Render("► LED: "))
		} else {
			sb.WriteString(inactiveControlStyle.Render("  LED: "))
		}
		sb.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.rgbwIndex+1)))
		sb.WriteString("\n")

		// R
		if m.activeControl == 3 && isActive {
			sb.WriteString(activeControlStyle.Render("► R: "))
		} else {
			sb.WriteString(inactiveControlStyle.Render("  R: "))
		}
		sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.r)))
		sb.WriteString("\n")

		// G
		if m.activeControl == 4 && isActive {
			sb.WriteString(activeControlStyle.Render("► G: "))
		} else {
			sb.WriteString(inactiveControlStyle.Render("  G: "))
		}
		sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.g)))
		sb.WriteString("\n")

		// B
		if m.activeControl == 5 && isActive {
			sb.WriteString(activeControlStyle.Render("► B: "))
		} else {
			sb.WriteString(inactiveControlStyle.Render("  B: "))
		}
		sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.b)))
		sb.WriteString("\n")

		// W
		if m.activeControl == 6 && isActive {
			sb.WriteString(activeControlStyle.Render("► W: "))
		} else {
			sb.WriteString(inactiveControlStyle.Render("  W: "))
		}
		sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.w)))

	} else { // White mode
		// LED Index
		if m.activeControl == 2 && isActive {
			sb.WriteString(activeControlStyle.Render("► LED: "))
		} else {
			sb.WriteString(inactiveControlStyle.Render("  LED: "))
		}
		sb.WriteString(valueStyle.Render(fmt.Sprintf("%d", m.whiteIndex+1)))
		sb.WriteString("\n")

		// Brightness
		if m.activeControl == 3 && isActive {
			sb.WriteString(activeControlStyle.Render("► Brightness: "))
		} else {
			sb.WriteString(inactiveControlStyle.Render("  Brightness: "))
		}
		sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.whiteBrightness)))
	}

	return sb.String()
}
