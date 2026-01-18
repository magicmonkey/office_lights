package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Messages for internal events

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type publishSuccessMsg struct{}
type publishErrorMsg struct{ err error }

// tickMsg is sent periodically to refresh values from drivers
type tickMsg time.Time

// tick returns a command that sends a tickMsg after a delay
func tick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
