package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/key"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	keys := DefaultKeyMap()

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, keys.NextSection):
			m.activeSection = (m.activeSection + 1) % SectionCount
			return m, nil

		case key.Matches(msg, keys.PrevSection):
			m.activeSection = (m.activeSection - 1 + SectionCount) % SectionCount
			return m, nil

		case key.Matches(msg, keys.Left):
			cmd = m.handleLeft()
			return m, cmd

		case key.Matches(msg, keys.Right):
			cmd = m.handleRight()
			return m, cmd

		case key.Matches(msg, keys.Up):
			cmd = m.handleAdjust(1)
			return m, cmd

		case key.Matches(msg, keys.Down):
			cmd = m.handleAdjust(-1)
			return m, cmd

		case key.Matches(msg, keys.BigUp):
			cmd = m.handleAdjust(10)
			return m, cmd

		case key.Matches(msg, keys.BigDown):
			cmd = m.handleAdjust(-10)
			return m, cmd

		case key.Matches(msg, keys.Toggle):
			cmd = m.handleToggle()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case publishSuccessMsg:
		// Successfully published
		m.err = nil
		return m, nil

	case publishErrorMsg:
		m.err = msg.err
		return m, nil

	case tickMsg:
		// Refresh all values from drivers
		m.ledStrip.refresh()
		m.ledBar.refresh()
		m.videoLight1.refresh()
		m.videoLight2.refresh()
		// Continue the tick cycle
		return m, tick()
	}

	return m, nil
}

func (m *Model) handleLeft() tea.Cmd {
	switch m.activeSection {
	case SectionLEDStrip:
		m.ledStrip.prevControl()
	case SectionLEDBar:
		m.ledBar.prevControl()
	case SectionVideoLight1:
		m.videoLight1.prevControl()
	case SectionVideoLight2:
		m.videoLight2.prevControl()
	}
	return nil
}

func (m *Model) handleRight() tea.Cmd {
	switch m.activeSection {
	case SectionLEDStrip:
		m.ledStrip.nextControl()
	case SectionLEDBar:
		m.ledBar.nextControl()
	case SectionVideoLight1:
		m.videoLight1.nextControl()
	case SectionVideoLight2:
		m.videoLight2.nextControl()
	}
	return nil
}

func (m *Model) handleAdjust(delta int) tea.Cmd {
	switch m.activeSection {
	case SectionLEDStrip:
		return m.ledStrip.adjustValue(delta)
	case SectionLEDBar:
		return m.ledBar.adjustValue(delta)
	case SectionVideoLight1:
		return m.videoLight1.adjustValue(delta)
	case SectionVideoLight2:
		return m.videoLight2.adjustValue(delta)
	}
	return nil
}

func (m *Model) handleToggle() tea.Cmd {
	switch m.activeSection {
	case SectionVideoLight1:
		return m.videoLight1.toggle()
	case SectionVideoLight2:
		return m.videoLight2.toggle()
	}
	return nil
}
