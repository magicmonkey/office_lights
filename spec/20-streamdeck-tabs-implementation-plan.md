# Stream Deck+ Tabs Implementation Plan

## Overview

This document provides step-by-step implementation tasks for adding tabbed navigation to the Stream Deck+ interface.

**Prerequisites:** Stream Deck+ interface implemented (specs 17-18)

## Implementation Phases

### Phase 1: State Model Updates

**Goal:** Add tab tracking to the StreamDeckUI struct

**File:** `streamdeck/model.go`

**Tasks:**

1. Add Tab type and constants:
```go
type Tab int

const (
    TabLightControl Tab = iota
    TabFuture2
    TabFuture3
    TabFuture4
)

func (t Tab) String() string {
    switch t {
    case TabLightControl:
        return "Lights"
    case TabFuture2:
        return "Tab 2"
    case TabFuture3:
        return "Tab 3"
    case TabFuture4:
        return "Tab 4"
    default:
        return "Unknown"
    }
}
```

2. Add `currentTab` field to StreamDeckUI struct:
```go
type StreamDeckUI struct {
    // ... existing fields ...
    currentTab  Tab   // Currently selected tab
    currentMode Mode  // Mode within TabLightControl
}
```

3. Initialize `currentTab` to `TabLightControl` in constructor

**Validation:**
- Code compiles
- Default tab is TabLightControl on startup

---

### Phase 2: Top Row Button Handling

**Goal:** Make top row buttons select tabs instead of being reserved

**File:** `streamdeck/events.go`

**Tasks:**

1. Update `handleButtonPress` to handle tab selection:
```go
func (s *StreamDeckUI) handleButtonPress(buttonIndex int) {
    log.Printf("Button %d pressed", buttonIndex)

    // Top row (0-3): Tab selection
    if buttonIndex < 4 {
        newTab := Tab(buttonIndex)
        if newTab != s.currentTab {
            log.Printf("Switching tab from %s to %s", s.currentTab, newTab)
            s.currentTab = newTab
            // Update all displays
            if err := s.updateButtons(); err != nil {
                log.Printf("Error updating buttons: %v", err)
            }
            if err := s.updateTouchscreen(); err != nil {
                log.Printf("Error updating touchscreen: %v", err)
            }
        }
        return
    }

    // Second row (4-7): Tab-specific actions
    switch s.currentTab {
    case TabLightControl:
        // Mode selection (existing behavior)
        newMode := Mode(buttonIndex - 4)
        if newMode != s.currentMode {
            log.Printf("Switching mode from %s to %s", s.currentMode, newMode)
            s.currentMode = newMode
            if err := s.updateButtons(); err != nil {
                log.Printf("Error updating buttons: %v", err)
            }
            if err := s.updateTouchscreen(); err != nil {
                log.Printf("Error updating touchscreen: %v", err)
            }
        }
    default:
        // Future tabs: no action yet
        log.Printf("Button %d pressed on unimplemented tab", buttonIndex)
    }
}
```

**Validation:**
- Pressing buttons 0-3 changes the active tab
- Tab changes are logged
- Display updates when tab changes

---

### Phase 3: Top Row Button Rendering

**Goal:** Render tab selection buttons on the top row

**File:** `streamdeck/render.go`

**Tasks:**

1. Update `renderButton` to render tab buttons on top row:
```go
func (s *StreamDeckUI) renderButton(index int) (image.Image, error) {
    // Top row (0-3): Tab buttons
    if index < 4 {
        return s.renderTabButton(index)
    }

    // Second row (4-7): Tab-specific buttons
    switch s.currentTab {
    case TabLightControl:
        return s.renderModeButton(index - 4)
    default:
        return s.renderBlankButton(), nil
    }
}
```

2. Add `renderTabButton` function:
```go
func (s *StreamDeckUI) renderTabButton(index int) (image.Image, error) {
    tab := Tab(index)
    isActive := s.currentTab == tab

    // Try to load tab icon
    iconPath := filepath.Join("streamdeck", "icons", s.getTabIconFilename(tab))
    img, err := loadImage(iconPath)
    if err != nil {
        // Fallback to text button
        return s.renderTextButton(tab.String(), isActive), nil
    }

    if isActive {
        img = applyHighlight(img)
    } else {
        img = applyDim(img)
    }

    return img, nil
}

func (s *StreamDeckUI) getTabIconFilename(tab Tab) string {
    switch tab {
    case TabLightControl:
        return "tab_lights.png"
    case TabFuture2:
        return "tab_2.png"
    case TabFuture3:
        return "tab_3.png"
    case TabFuture4:
        return "tab_4.png"
    default:
        return "unknown.png"
    }
}
```

**Validation:**
- Top row shows tab buttons
- Active tab is highlighted
- Inactive tabs are dimmed

---

### Phase 4: Second Row Button Rendering

**Goal:** Render appropriate buttons based on active tab

**File:** `streamdeck/render.go`

**Tasks:**

1. Update second row rendering (already in renderButton from Phase 3):
   - Tab 1: Show mode selection buttons (existing behavior)
   - Tabs 2-4: Show blank buttons

**Validation:**
- Tab 1: Second row shows LED Strip, LED Bar RGBW, LED Bar White, Video Lights
- Tabs 2-4: Second row shows blank buttons

---

### Phase 5: Touchscreen Rendering

**Goal:** Render tab-specific content on the touchscreen

**File:** `streamdeck/render.go`

**Tasks:**

1. Update `renderTouchscreen` to be tab-aware:
```go
func (s *StreamDeckUI) renderTouchscreen() image.Image {
    switch s.currentTab {
    case TabLightControl:
        return s.renderLightControlTouchscreen()
    default:
        return s.renderPlaceholderTouchscreen()
    }
}

func (s *StreamDeckUI) renderLightControlTouchscreen() image.Image {
    // Existing touchscreen rendering logic
    img := image.NewRGBA(image.Rect(0, 0, touchWidth, touchHeight))
    draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{20, 20, 20, 255}}, image.Point{}, draw.Src)

    sections := s.getSectionData()
    for i := 0; i < 4; i++ {
        s.renderSection(img, i, sections[i])
    }

    return img
}

func (s *StreamDeckUI) renderPlaceholderTouchscreen() image.Image {
    img := image.NewRGBA(image.Rect(0, 0, touchWidth, touchHeight))

    // Dark background
    draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{30, 30, 30, 255}}, image.Point{}, draw.Src)

    // Center text "Coming Soon"
    text := "Coming Soon"
    drawTextAt(img, text, touchWidth/2, touchHeight/2, color.RGBA{100, 100, 100, 255}, true)

    return img
}
```

**Validation:**
- Tab 1: Touchscreen shows light control sections
- Tabs 2-4: Touchscreen shows "Coming Soon" placeholder

---

### Phase 6: Dial Handling

**Goal:** Make dials tab-aware

**File:** `streamdeck/events.go`

**Tasks:**

1. Update `handleDialRotate` to check current tab:
```go
func (s *StreamDeckUI) handleDialRotate(dialIndex int, ticks int) {
    if dialIndex < 0 || dialIndex > 3 {
        log.Printf("Invalid dial index: %d", dialIndex)
        return
    }

    // Only handle dials on Tab 1
    if s.currentTab != TabLightControl {
        log.Printf("Dial %d rotated on unimplemented tab", dialIndex)
        return
    }

    // Existing dial rotation logic...
    log.Printf("Dial %d rotated: %d ticks", dialIndex, ticks)
    // ... rest of existing code
}
```

2. Update `handleDialPress` similarly:
```go
func (s *StreamDeckUI) handleDialPress(dialIndex int) {
    if dialIndex < 0 || dialIndex > 3 {
        log.Printf("Invalid dial index: %d", dialIndex)
        return
    }

    // Only handle dial presses on Tab 1
    if s.currentTab != TabLightControl {
        log.Printf("Dial %d pressed on unimplemented tab", dialIndex)
        return
    }

    // Existing dial press logic...
    log.Printf("Dial %d pressed", dialIndex)
    // ... rest of existing code
}
```

**Validation:**
- Tab 1: Dials adjust light values (existing behavior)
- Tabs 2-4: Dial actions are logged but have no effect

---

### Phase 7: Tab Icons (Optional)

**Goal:** Create placeholder icons for tabs

**Directory:** `streamdeck/icons/`

**Tasks:**

1. Create or source 120x120 PNG icons:
   - `tab_lights.png` - Light bulb or sun icon for Tab 1
   - `tab_2.png` - Placeholder for Tab 2
   - `tab_3.png` - Placeholder for Tab 3
   - `tab_4.png` - Placeholder for Tab 4

2. Alternatively, use text-based buttons (already implemented as fallback)

**Validation:**
- Tab icons display correctly (or text fallback works)

---

### Phase 8: Testing

**Goal:** Verify all tab functionality

**Test Cases:**

1. **Tab Selection:**
   - [ ] Button 0 selects Tab 1 (Light Control)
   - [ ] Button 1 selects Tab 2 (placeholder)
   - [ ] Button 2 selects Tab 3 (placeholder)
   - [ ] Button 3 selects Tab 4 (placeholder)
   - [ ] Active tab button is highlighted
   - [ ] Inactive tab buttons are dimmed

2. **Tab 1 Functionality:**
   - [ ] Second row shows mode selection buttons
   - [ ] Mode selection works as before
   - [ ] Touchscreen shows light values
   - [ ] Dials adjust values
   - [ ] Dial clicks toggle values

3. **Tabs 2-4 Placeholder:**
   - [ ] Second row shows blank buttons
   - [ ] Touchscreen shows "Coming Soon"
   - [ ] Dials have no effect
   - [ ] Second row buttons have no effect

4. **Tab Switching:**
   - [ ] Switching tabs updates all displays
   - [ ] Returning to Tab 1 restores previous mode
   - [ ] No state is lost when switching tabs

**Validation:**
- All test cases pass
- No regressions in existing functionality

---

## Summary of File Changes

| File | Changes |
|------|---------|
| `streamdeck/model.go` | Add Tab type, constants, String() method; add currentTab field |
| `streamdeck/events.go` | Update handleButtonPress for tabs; add tab checks to dial handlers |
| `streamdeck/render.go` | Add renderTabButton; update renderButton; add renderPlaceholderTouchscreen |
| `streamdeck/icons/` | Add tab icon files (optional) |

## Implementation Order

1. Phase 1: State model updates (model.go)
2. Phase 2: Button event handling (events.go)
3. Phase 3: Top row rendering (render.go)
4. Phase 4: Second row rendering (render.go)
5. Phase 5: Touchscreen rendering (render.go)
6. Phase 6: Dial handling (events.go)
7. Phase 7: Tab icons (optional)
8. Phase 8: Testing

## Notes

- The existing light control functionality (Tab 1) should remain unchanged
- Tabs 2-4 are placeholders for future features
- When implementing future tabs, follow the same pattern: update events, update rendering
- Consider adding tab state persistence to the database (future enhancement)
