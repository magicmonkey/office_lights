# TUI Architecture

## Overview

Implement a terminal-based user interface (TUI) that provides interactive control for all office lights. The interface displays all lights simultaneously in a split-screen layout with keyboard-based navigation and control.

## Requirements

### Layout Structure

The screen is divided into 4 sections:

1. **LED Strip Section** (Top Left)
   - RGB color controls (0-255 for each channel)
   - Red slider/input
   - Green slider/input
   - Blue slider/input

2. **LED Bar Section** (Top Right)
   - Section selector (Section 1 or Section 2)
   - RGBW controls for selected LED (1-6 in that section)
   - LED index selector
   - Red, Green, Blue, White sliders (0-255)
   - White LED controls (1-13 in that section)
   - White LED index selector
   - Brightness slider (0-255)

3. **Video Light 1 Section** (Bottom Left)
   - On/Off toggle
   - Brightness slider (0-100)

4. **Video Light 2 Section** (Bottom Right)
   - On/Off toggle
   - Brightness slider (0-100)

### Navigation

- **TAB key**: Switch focus between the 4 main sections (LED Strip → LED Bar → Video Light 1 → Video Light 2 → LED Strip...)
- **Shift+TAB**: Switch focus in reverse order
- **Arrow keys** (within a section):
  - Left/Right: Move between different controls
  - Up/Down: Increment/decrement values by small amounts (e.g., ±1)
  - Shift+Up/Down: Increment/decrement values by large amounts (e.g., ±10)
- **Enter**: Toggle boolean values (on/off states)
- **ESC** or **Ctrl+C**: Exit the application

### Value Adjustment

- **Small increment**: ±1 (Up/Down arrows)
- **Large increment**: ±10 (Shift+Up/Down arrows)
- **Bounds checking**: All values stay within valid ranges
  - RGB/RGBW: 0-255
  - Video light brightness: 0-100
  - LED indices: Within valid ranges (1-6 for RGBW, 1-13 for white)

### Real-time Updates

- When a value is changed, immediately:
  1. Update the internal driver state
  2. Publish to MQTT
  3. Save to database
  4. Update the TUI display

### Visual Feedback

- **Active section**: Highlighted border or different color
- **Active control**: Cursor/highlight indicator
- **Current values**: Displayed numerically and/or as progress bars
- **Status indicators**: Show connection status, last update time

## Technology Selection

### Recommended Library: Bubbletea

Use the **Bubbletea** framework (`github.com/charmbracelet/bubbletea`) for the following reasons:

1. **Modern Go TUI framework** - Designed specifically for Go
2. **Elm Architecture** - Clean separation of model, update, and view
3. **Composable** - Easy to create separate components for each section
4. **Active development** - Well-maintained with good documentation
5. **Rich ecosystem** - Works with Bubbles (component library) and Lipgloss (styling)

**Supporting libraries:**
- `github.com/charmbracelet/bubbles` - Pre-built components (progress bars, text inputs, etc.)
- `github.com/charmbracelet/lipgloss` - Styling and layout

### Alternative Libraries (Not Recommended)

- **tview** - More widget-focused, less flexible for custom layouts
- **termui** - Lower-level, more manual work required
- **tcell** - Too low-level, would need to build everything from scratch

## Architecture Design

### Package Structure

```
tui/
├── tui.go           # Main TUI orchestrator, Bubbletea program
├── model.go         # Root model (Bubbletea model interface)
├── ledstrip.go      # LED strip section component
├── ledbar.go        # LED bar section component
├── videolight.go    # Video light section component
├── layout.go        # Layout and rendering logic
└── keys.go          # Key binding definitions
```

### Component Hierarchy

```
Root Model (tui.Model)
├── LED Strip Component
│   ├── R value control
│   ├── G value control
│   └── B value control
├── LED Bar Component
│   ├── Section selector (1 or 2)
│   ├── RGBW controls
│   │   ├── LED index selector
│   │   ├── R value control
│   │   ├── G value control
│   │   ├── B value control
│   │   └── W value control
│   └── White LED controls
│       ├── LED index selector
│       └── Brightness control
├── Video Light 1 Component
│   ├── On/Off toggle
│   └── Brightness control
└── Video Light 2 Component
    ├── On/Off toggle
    └── Brightness control
```

### Bubbletea Model Structure

```go
type Model struct {
    // Focus management
    activeSection int  // 0=LED Strip, 1=LED Bar, 2=VL1, 3=VL2

    // Component models
    ledStrip    *ledStripModel
    ledBar      *ledBarModel
    videoLight1 *videoLightModel
    videoLight2 *videoLightModel

    // Driver references (for publishing changes)
    stripDriver *ledstrip.LEDStrip
    barDriver   *ledbar.LEDBar
    vl1Driver   *videolight.VideoLight
    vl2Driver   *videolight.VideoLight

    // UI state
    width  int
    height int
    ready  bool
}
```

### Message Passing

Bubbletea uses a message-passing architecture:

```go
// Messages for internal events
type FocusChangedMsg int
type ValueChangedMsg struct {
    Section string
    Control string
    Value   int
}

// Messages for driver updates
type PublishSuccessMsg struct{}
type PublishErrorMsg struct{ Err error }
```

### Update Flow

1. **User Input** → Key press generates Bubbletea `tea.KeyMsg`
2. **Update Function** → Processes key, updates model, triggers driver action
3. **Driver Action** → Publishes to MQTT, saves to database
4. **View Function** → Re-renders the TUI with updated values

## Integration with Existing Code

### Main Application Changes

The TUI will be an optional mode:

```go
func main() {
    // ... existing setup ...

    // Check if TUI mode is requested
    if os.Getenv("TUI") != "" || len(os.Args) > 1 && os.Args[1] == "tui" {
        runTUI(ledStrip, ledBar, videoLight1, videoLight2)
        return
    }

    // ... existing main loop ...
}
```

### Driver Integration

The TUI will use existing driver methods:
- `ledStrip.SetColor(r, g, b)` - Updates and publishes
- `ledBar.SetRGBW(section, index, r, g, b, w)` - Updates and publishes
- `ledBar.SetWhite(section, index, brightness)` - Updates and publishes
- `videoLight.TurnOn(brightness)` / `TurnOff()` - Updates and publishes

**Important**: The TUI does NOT need to directly handle MQTT or database operations - these are already handled by the driver methods.

## Error Handling

- **MQTT Connection Lost**: Display warning banner, queue changes
- **Database Write Failure**: Log error, continue operation (MQTT still works)
- **Invalid Input**: Silently clamp to valid range, don't allow invalid values
- **Terminal Resize**: Handle `tea.WindowSizeMsg` to adjust layout

## Testing Strategy

### Unit Tests

- Test each component model independently
- Test update logic for all key combinations
- Test value clamping and validation
- Mock driver interfaces for testing

### Integration Tests

- Test with mock MQTT publisher
- Verify driver methods are called correctly
- Test state synchronization

### Manual Testing

- Keyboard navigation through all sections
- Value changes at boundaries (0, 255, etc.)
- Rapid key presses (debouncing)
- Terminal resize behavior
- Exit and cleanup

## Performance Considerations

- **Rendering**: Only re-render when state changes
- **Debouncing**: Avoid publishing on every keystroke if value changes rapidly
- **Efficient Updates**: Use Bubbletea's efficient rendering (only diffs)

## Accessibility

- Clear visual indicators for focused section and control
- Consistent keyboard navigation patterns
- Help text showing available keys
- High contrast colors for readability

## Future Enhancements (Not in Initial Implementation)

- Mouse support for clicking controls
- Color preview (background color matching RGB values)
- Preset scenes (save/load combinations)
- History/undo functionality
- Real-time preview animation
