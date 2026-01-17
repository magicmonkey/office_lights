package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary   = lipgloss.Color("#00ADD8") // Go blue
	colorSecondary = lipgloss.Color("#FDDD00") // Go yellow
	colorActive    = lipgloss.Color("#00FF00") // Green for active
	colorInactive  = lipgloss.Color("#888888") // Gray for inactive
	colorBorder    = lipgloss.Color("#444444")

	// Base styles
	baseStyle = lipgloss.NewStyle().
		Padding(0, 1)

	// Section styles
	activeSectionStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorActive).
		Padding(1, 2)

	inactiveSectionStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Padding(1, 2)

	// Title styles
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		MarginBottom(1)

	// Control styles
	activeControlStyle = lipgloss.NewStyle().
		Foreground(colorActive).
		Bold(true)

	inactiveControlStyle = lipgloss.NewStyle().
		Foreground(colorInactive)

	// Value styles
	valueStyle = lipgloss.NewStyle().
		Foreground(colorSecondary).
		Bold(true)

	// Help text style
	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Italic(true)
)
