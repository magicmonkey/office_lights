# Stream Deck+ Tabs Architecture

## Overview

This specification extends the Stream Deck+ interface to support a tabbed navigation system. The top row of buttons (0-3) now selects between 4 different "tabs" or "pages", each providing distinct functionality.

## Tab System Design

### Button Layout

```
Top Row (Tab Selection):
┌──────────┬──────────┬──────────┬──────────┐
│  Tab 1   │  Tab 2   │  Tab 3   │  Tab 4   │
│ (Active) │ (Future) │ (Future) │ (Future) │
└──────────┴──────────┴──────────┴──────────┘

Second Row (Tab-Specific Controls):
┌──────────┬──────────┬──────────┬──────────┐
│ Control1 │ Control2 │ Control3 │ Control4 │
│  (varies by tab)                          │
└──────────┴──────────┴──────────┴──────────┘
```

### Tab Definitions

| Tab | Button | Name | Status | Description |
|-----|--------|------|--------|-------------|
| 1 | 0 | Light Control | Implemented | Individual light control with mode selection |
| 2 | 1 | (Undefined) | Future | Reserved for future functionality |
| 3 | 2 | (Undefined) | Future | Reserved for future functionality |
| 4 | 3 | (Undefined) | Future | Reserved for future functionality |

## Tab 1: Light Control

This is the existing functionality, now contained within Tab 1.

### Second Row Buttons (Mode Selection)

When Tab 1 is active, the second row buttons act as a radio button set:

| Button 4 | Button 5 | Button 6 | Button 7 |
|----------|----------|----------|----------|
| LED Strip | LED Bar RGBW | LED Bar White | Video Lights |

### Touchscreen Display

The touchscreen shows 4 sections based on the selected mode:

- **LED Strip Mode**: Red, Green, Blue, (empty)
- **LED Bar RGBW Mode**: Red, Green, Blue, White
- **LED Bar White Mode**: Section 1, Section 2, (empty), (empty)
- **Video Lights Mode**: Light 1, Light 2, (empty), (empty)

### Dial Behavior

- Rotation: Adjust values in increments of 5
- Click: Toggle between 0 and last value (or on/off for video lights)

## Tabs 2-4: Future Functionality

Reserved for future features. Possible uses include:

- **Presets/Scenes**: Quick-access lighting presets
- **Timers**: Scheduled lighting changes
- **Effects**: Animated lighting effects
- **Settings**: Configuration options

When selected, these tabs should display a placeholder indicating they are not yet implemented.

## Technical Architecture

### State Model Changes

```go
type Tab int

const (
    TabLightControl Tab = iota  // Tab 1: Existing light control
    TabFuture2                   // Tab 2: Reserved
    TabFuture3                   // Tab 3: Reserved
    TabFuture4                   // Tab 4: Reserved
)

type StreamDeckUI struct {
    // ... existing fields ...

    currentTab  Tab   // Currently selected tab (0-3)
    currentMode Mode  // Mode within Tab 1 (LED Strip, LED Bar RGBW, etc.)
}
```

### Event Handling Changes

**Top Row Buttons (0-3):**
- Now select the active tab
- Update button visuals to show active tab
- Re-render second row buttons based on selected tab
- Re-render touchscreen based on selected tab

**Second Row Buttons (4-7):**
- Behavior depends on the active tab
- Tab 1: Mode selection (existing behavior)
- Tabs 2-4: No action (or tab-specific actions when implemented)

**Dials:**
- Behavior depends on the active tab
- Tab 1: Adjust light values (existing behavior)
- Tabs 2-4: No action (or tab-specific actions when implemented)

### Rendering Changes

**Top Row Button Rendering:**
```go
func (s *StreamDeckUI) renderTabButton(tabIndex int) image.Image {
    isActive := s.currentTab == Tab(tabIndex)
    // Render with active/inactive styling
    // Show tab name or icon
}
```

**Second Row Button Rendering:**
```go
func (s *StreamDeckUI) renderSecondRowButton(buttonIndex int) image.Image {
    switch s.currentTab {
    case TabLightControl:
        return s.renderModeButton(buttonIndex) // Existing behavior
    default:
        return s.renderBlankButton() // Placeholder for future tabs
    }
}
```

**Touchscreen Rendering:**
```go
func (s *StreamDeckUI) renderTouchscreen() image.Image {
    switch s.currentTab {
    case TabLightControl:
        return s.renderLightControlTouchscreen() // Existing behavior
    default:
        return s.renderPlaceholderTouchscreen("Coming Soon")
    }
}
```

## Visual Design

### Tab Button States

- **Active Tab**: Bright/highlighted, distinct background color
- **Inactive Tab**: Dimmed, darker background
- **Future Tabs**: Show icon or text indicating "coming soon" or placeholder

### Tab Icons (Suggested)

| Tab | Icon Suggestion |
|-----|-----------------|
| 1 | Light bulb or sun icon |
| 2 | Clock or timer icon |
| 3 | Sparkle or effect icon |
| 4 | Gear or settings icon |

## Implementation Considerations

### Backward Compatibility

The existing light control functionality remains unchanged within Tab 1. The main changes are:

1. Top row buttons now select tabs instead of being reserved
2. Tab 1 is selected by default on startup
3. All existing mode selection moves to second row (already done)

### Default State

- On startup, Tab 1 (Light Control) is selected
- Within Tab 1, LED Strip mode is selected by default (existing behavior)

### Persistence

- Consider saving the last-selected tab to the database
- On startup, restore to the last-used tab (optional enhancement)

## Error Handling

- If a user selects an unimplemented tab, show a placeholder on the touchscreen
- Second row buttons should be blank or show "N/A" for unimplemented tabs
- Dials should have no effect on unimplemented tabs

## Testing

### Unit Tests

- Tab switching updates `currentTab` correctly
- Top row buttons trigger tab changes
- Second row buttons behave correctly per tab
- Touchscreen renders correct content per tab

### Integration Tests

- Tab 1 functionality unchanged from existing behavior
- Tabs 2-4 show placeholder content
- Tab state persists correctly (if implemented)

### Manual Tests

- [ ] Button 0 selects Tab 1 (Light Control)
- [ ] Buttons 1-3 select placeholder tabs
- [ ] Active tab button is visually highlighted
- [ ] Second row shows mode buttons only on Tab 1
- [ ] Touchscreen shows light controls only on Tab 1
- [ ] Touchscreen shows "Coming Soon" on Tabs 2-4
- [ ] Dials work on Tab 1, no effect on Tabs 2-4

## Future Extensibility

When implementing Tabs 2-4, the pattern will be:

1. Define the tab's purpose and controls
2. Add tab-specific rendering for second row buttons
3. Add tab-specific touchscreen content
4. Add tab-specific dial behavior
5. Update event handlers for the new tab

The architecture is designed to make adding new tabs straightforward without modifying the core tab-switching logic.
