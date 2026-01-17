package tui

import (
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.ready {
		return "Initializing TUI..."
	}

	// Render each section
	ledStripView := m.renderSection(
		m.ledStrip.View(m.activeSection == SectionLEDStrip),
		m.activeSection == SectionLEDStrip,
	)

	ledBarView := m.renderSection(
		m.ledBar.View(m.activeSection == SectionLEDBar),
		m.activeSection == SectionLEDBar,
	)

	vl1View := m.renderSection(
		m.videoLight1.View(m.activeSection == SectionVideoLight1),
		m.activeSection == SectionVideoLight1,
	)

	vl2View := m.renderSection(
		m.videoLight2.View(m.activeSection == SectionVideoLight2),
		m.activeSection == SectionVideoLight2,
	)

	// Layout: 2x2 grid
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, ledStripView, ledBarView)
	bottomRow := lipgloss.JoinHorizontal(lipgloss.Top, vl1View, vl2View)
	content := lipgloss.JoinVertical(lipgloss.Left, topRow, bottomRow)

	// Add help text at bottom
	help := m.renderHelp()

	return lipgloss.JoinVertical(lipgloss.Left, content, help)
}

func (m Model) renderSection(content string, isActive bool) string {
	style := inactiveSectionStyle
	if isActive {
		style = activeSectionStyle
	}

	// Calculate section dimensions (half screen width/height)
	width := (m.width / 2) - 6
	height := (m.height / 2) - 5

	// Ensure minimum dimensions
	if width < 20 {
		width = 20
	}
	if height < 8 {
		height = 8
	}

	return style.Width(width).Height(height).Render(content)
}

func (m Model) renderHelp() string {
	help := "TAB: next section | ←→: select control | ↑↓: adjust (+1) | Shift+↑↓: adjust (+10) | Enter: toggle | ESC: quit"
	if m.err != nil {
		help = "Error: " + m.err.Error() + " | " + help
	}
	return "\n" + helpStyle.Render(help)
}
