# TUI Implementation Plan

## Overview

Step-by-step guide for implementing the terminal user interface (TUI) for the office lights control system.

## Prerequisites

- Completed phases 1-10 (all drivers, MQTT, and state storage)
- Go 1.24.0 or later
- Terminal with ANSI color support

## Implementation Phases

### Phase 1: Dependencies and Project Structure

**Goal**: Install required libraries and create TUI package structure

#### Step 1.1: Install Dependencies

```bash
go get github.com/charmbracelet/bubbletea@latest
go get github.com/charmbracelet/bubbles@latest
go get github.com/charmbracelet/lipgloss@latest
```

#### Step 1.2: Create TUI Package Structure

Create the following files:

```
tui/
├── tui.go           # Main entry point and program initialization
├── model.go         # Root Bubbletea model
├── update.go        # Update function (message handling)
├── view.go          # View function (rendering)
├── keys.go          # Key binding definitions
├── ledstrip.go      # LED strip component
├── ledbar.go        # LED bar component
├── videolight.go    # Video light component
├── styles.go        # Lipgloss styles and theming
└── messages.go      # Custom message types
```

#### Verification

- All dependencies installed without errors
- Package structure created
- Files compile (can have empty/stub implementations)

---

### Phase 2: Core TUI Model and Key Bindings

**Goal**: Implement the root model and basic key handling

#### Step 2.1: Define Model Structure (`model.go`)

```go
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
    return nil
}
```

#### Step 2.2: Define Key Bindings (`keys.go`)

```go
package tui

import "github.com/charmbracelet/bubbles/key"

// KeyMap defines all key bindings
type KeyMap struct {
    NextSection key.Binding
    PrevSection key.Binding
    Up          key.Binding
    Down        key.Binding
    Left        key.Binding
    Right       key.Binding
    BigUp       key.Binding
    BigDown     key.Binding
    Toggle      key.Binding
    Quit        key.Binding
}

// DefaultKeyMap returns the default key bindings
func DefaultKeyMap() KeyMap {
    return KeyMap{
        NextSection: key.NewBinding(
            key.WithKeys("tab"),
            key.WithHelp("tab", "next section"),
        ),
        PrevSection: key.NewBinding(
            key.WithKeys("shift+tab"),
            key.WithHelp("shift+tab", "prev section"),
        ),
        Up: key.NewBinding(
            key.WithKeys("up"),
            key.WithHelp("↑", "increase (+1)"),
        ),
        Down: key.NewBinding(
            key.WithKeys("down"),
            key.WithHelp("↓", "decrease (-1)"),
        ),
        Left: key.NewBinding(
            key.WithKeys("left"),
            key.WithHelp("←", "prev control"),
        ),
        Right: key.NewBinding(
            key.WithKeys("right"),
            key.WithHelp("→", "next control"),
        ),
        BigUp: key.NewBinding(
            key.WithKeys("shift+up"),
            key.WithHelp("shift+↑", "increase (+10)"),
        ),
        BigDown: key.NewBinding(
            key.WithKeys("shift+down"),
            key.WithHelp("shift+↓", "decrease (-10)"),
        ),
        Toggle: key.NewBinding(
            key.WithKeys("enter"),
            key.WithHelp("enter", "toggle on/off"),
        ),
        Quit: key.NewBinding(
            key.WithKeys("esc", "ctrl+c"),
            key.WithHelp("esc", "quit"),
        ),
    }
}
```

#### Step 2.3: Define Custom Messages (`messages.go`)

```go
package tui

// Messages for internal events

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type publishSuccessMsg struct{}
type publishErrorMsg struct{ err error }
```

#### Verification

- Model compiles
- Key bindings defined
- Messages defined
- Can create a new Model instance

---

### Phase 3: Component Models

**Goal**: Implement individual component models for each light type

#### Step 3.1: LED Strip Component (`ledstrip.go`)

```go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kevin/office_lights/drivers/ledstrip"
)

// ledStripModel represents the LED strip section
type ledStripModel struct {
    driver       *ledstrip.LEDStrip
    activeControl int // 0=R, 1=G, 2=B
    r, g, b      int
}

func newLEDStripModel(driver *ledstrip.LEDStrip) ledStripModel {
    // Load current state from driver
    return ledStripModel{
        driver:        driver,
        activeControl: 0,
        r:             driver.R(), // Assuming getter methods exist
        g:             driver.G(),
        b:             driver.B(),
    }
}

// Update handles messages for this component
func (m ledStripModel) Update(msg tea.Msg) (ledStripModel, tea.Cmd) {
    // Will be implemented in update.go
    return m, nil
}

// View renders this component
func (m ledStripModel) View() string {
    // Will be implemented in view.go
    return ""
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

func clamp(val, min, max int) int {
    if val < min {
        return min
    }
    if val > max {
        return max
    }
    return val
}
```

**Note**: This assumes driver methods exist to get current values. If not, we'll need to add getter methods to drivers.

#### Step 3.2: Video Light Component (`videolight.go`)

```go
package tui

import (
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
    return videoLightModel{
        driver:        driver,
        lightID:       lightID,
        activeControl: 0,
        on:            driver.IsOn(), // Assuming getter exists
        brightness:    driver.Brightness(),
    }
}

func (m *videoLightModel) nextControl() {
    m.activeControl = (m.activeControl + 1) % 2
}

func (m *videoLightModel) prevControl() {
    m.activeControl = (m.activeControl - 1 + 2) % 2
}

func (m *videoLightModel) adjustValue(delta int) tea.Cmd {
    if m.activeControl == 0 {
        // On/off toggle - ignore delta
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
```

#### Step 3.3: LED Bar Component (`ledbar.go`)

```go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/kevin/office_lights/drivers/ledbar"
)

// ledBarModel represents the LED bar section
type ledBarModel struct {
    driver        *ledbar.LEDBar
    mode          int // 0=RGBW, 1=White
    section       int // 1 or 2
    activeControl int // Depends on mode

    // RGBW mode
    rgbwIndex int // 0-5 (LED 1-6)
    r, g, b, w int

    // White mode
    whiteIndex     int // 0-12 (LED 1-13)
    whiteBrightness int
}

func newLEDBarModel(driver *ledbar.LEDBar) ledBarModel {
    // Initialize with default values
    // Load current state from driver if possible
    return ledBarModel{
        driver:  driver,
        mode:    0, // Start in RGBW mode
        section: 1, // Start with section 1
    }
}

func (m *ledBarModel) nextControl() {
    if m.mode == 0 { // RGBW mode
        m.activeControl = (m.activeControl + 1) % 5 // section, index, R, G, B, W
    } else { // White mode
        m.activeControl = (m.activeControl + 1) % 3 // section, index, brightness
    }
}

func (m *ledBarModel) prevControl() {
    if m.mode == 0 {
        m.activeControl = (m.activeControl - 1 + 5) % 5
    } else {
        m.activeControl = (m.activeControl - 1 + 3) % 3
    }
}

func (m *ledBarModel) adjustValue(delta int) tea.Cmd {
    // Implementation depends on active control
    // Update appropriate value and publish
    return m.publish()
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
```

#### Verification

- All component models compile
- Can create instances of each component
- Helper methods work correctly

---

### Phase 4: Update Function

**Goal**: Implement the root update function that handles all messages

#### Step 4.1: Implement Root Update (`update.go`)

```go
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
        return m, nil

    case publishErrorMsg:
        m.err = msg.err
        return m, nil
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
```

#### Verification

- Update function compiles
- All key bindings route to correct handlers
- Section switching works
- Component methods are called correctly

---

### Phase 5: View Functions and Styling

**Goal**: Implement rendering logic with proper layout and styling

#### Step 5.1: Define Styles (`styles.go`)

```go
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
```

#### Step 5.2: Implement Component Views

Each component (ledstrip.go, ledbar.go, videolight.go) needs a `View()` method:

```go
// In ledstrip.go
func (m ledStripModel) View(isActive bool) string {
    var sb strings.Builder

    sb.WriteString(titleStyle.Render("LED Strip"))
    sb.WriteString("\n\n")

    // Red control
    if m.activeControl == 0 {
        sb.WriteString(activeControlStyle.Render("► R: "))
    } else {
        sb.WriteString(inactiveControlStyle.Render("  R: "))
    }
    sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.r)))
    sb.WriteString("\n")

    // Green control
    if m.activeControl == 1 {
        sb.WriteString(activeControlStyle.Render("► G: "))
    } else {
        sb.WriteString(inactiveControlStyle.Render("  G: "))
    }
    sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.g)))
    sb.WriteString("\n")

    // Blue control
    if m.activeControl == 2 {
        sb.WriteString(activeControlStyle.Render("► B: "))
    } else {
        sb.WriteString(inactiveControlStyle.Render("  B: "))
    }
    sb.WriteString(valueStyle.Render(fmt.Sprintf("%3d", m.b)))

    return sb.String()
}
```

Similar implementations for video light and LED bar components.

#### Step 5.3: Implement Root View (`view.go`)

```go
package tui

import (
    "github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
    if !m.ready {
        return "Initializing..."
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
    width := (m.width / 2) - 4
    height := (m.height / 2) - 4

    return style.Width(width).Height(height).Render(content)
}

func (m Model) renderHelp() string {
    help := "TAB: next section | ←→: select control | ↑↓: adjust (+1) | Shift+↑↓: adjust (+10) | Enter: toggle | ESC: quit"
    if m.err != nil {
        help = "Error: " + m.err.Error() + " | " + help
    }
    return "\n" + helpStyle.Render(help)
}
```

#### Verification

- View renders correctly
- Layout is split into 4 sections
- Active section is highlighted
- Active control within section is highlighted
- Help text displays

---

### Phase 6: Main Integration and Entry Point

**Goal**: Create the TUI entry point and integrate with main.go

#### Step 6.1: Create TUI Entry Point (`tui.go`)

```go
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
```

#### Step 6.2: Update main.go

Add TUI mode to main.go:

```go
// After creating all drivers...

// Check if TUI mode is requested
useTUI := false
if len(os.Args) > 1 && os.Args[1] == "tui" {
    useTUI = true
}
if os.Getenv("TUI") != "" {
    useTUI = true
}

if useTUI {
    log.Println("Starting TUI mode...")
    if err := tui.Run(ledStrip, ledBar, videoLight1, videoLight2); err != nil {
        log.Fatalf("TUI error: %v", err)
    }
    return
}

// ... rest of main (existing behavior)
```

#### Verification

- Can run with `./office_lights tui`
- Can run with `TUI=1 ./office_lights`
- TUI starts and displays correctly
- Can exit with ESC or Ctrl+C

---

### Phase 7: Driver Getter Methods

**Goal**: Add getter methods to drivers if they don't exist

If drivers don't already have getter methods for current state, add them:

#### LED Strip (`drivers/ledstrip/ledstrip.go`)

```go
// R returns the current red value
func (l *LEDStrip) R() int {
    return l.r
}

// G returns the current green value
func (l *LEDStrip) G() int {
    return l.g
}

// B returns the current blue value
func (l *LEDStrip) B() int {
    return l.b
}
```

#### Video Light (`drivers/videolight/videolight.go`)

```go
// IsOn returns whether the light is on
func (v *VideoLight) IsOn() bool {
    return v.on
}

// Brightness returns the current brightness
func (v *VideoLight) Brightness() int {
    return v.brightness
}
```

#### LED Bar (`drivers/ledbar/ledbar.go`)

May need methods to get current RGBW and white values if not already present.

#### Verification

- All drivers have getter methods
- Getters return correct current state
- TUI can initialize with current state

---

### Phase 8: Testing

**Goal**: Test the TUI thoroughly

#### Unit Tests

Create `tui/tui_test.go`:

- Test key binding routing
- Test value clamping
- Test section switching
- Test component state updates

#### Integration Tests

- Start TUI with mock drivers
- Simulate key presses
- Verify driver methods are called

#### Manual Testing Checklist

- [ ] Navigation between all 4 sections with TAB
- [ ] Reverse navigation with Shift+TAB
- [ ] Left/Right arrows move between controls
- [ ] Up/Down adjust values by 1
- [ ] Shift+Up/Down adjust values by 10
- [ ] Values clamp at boundaries (0, 255, etc.)
- [ ] Enter toggles video light on/off
- [ ] ESC exits the application
- [ ] Ctrl+C exits the application
- [ ] Terminal resize handled gracefully
- [ ] All changes publish to MQTT
- [ ] All changes save to database
- [ ] Error messages display for MQTT failures

---

## Implementation Order Summary

1. **Phase 1**: Install dependencies, create package structure
2. **Phase 2**: Implement model, keys, messages
3. **Phase 3**: Implement component models
4. **Phase 4**: Implement update function
5. **Phase 5**: Implement view functions and styling
6. **Phase 6**: Create entry point and integrate with main
7. **Phase 7**: Add driver getter methods
8. **Phase 8**: Test thoroughly

## Success Criteria

- TUI runs without errors
- All 4 sections display correctly
- Keyboard navigation works as specified
- All value changes publish to MQTT
- All value changes save to database
- Terminal resize doesn't crash
- Clean exit on ESC/Ctrl+C
- Visual feedback is clear and helpful

## Estimated Complexity

- **Simple**: Phases 1-2 (setup and structure)
- **Moderate**: Phases 3-4 (components and update logic)
- **Moderate**: Phase 5 (view and styling)
- **Simple**: Phases 6-7 (integration)
- **Moderate**: Phase 8 (testing)

## Notes

- LED Bar component is the most complex due to multiple modes
- Consider implementing LED Strip and Video Light components first
- Test each component independently before integration
- Use `tea.WithAltScreen()` to preserve terminal state on exit
