# Stream Deck+ Scenes Implementation Plan (Tab 2)

## Overview

This document provides step-by-step implementation tasks for adding the Scenes feature (Tab 2) to the Stream Deck+ interface.

**Prerequisites:**
- Stream Deck+ tabs implemented (specs 19-20)
- Storage layer implemented (specs 9-12)

## Implementation Phases

### Phase 1: Database Schema Updates

**Goal:** Add scene tables to the SQLite database

**File:** `storage/schema.go`

**Tasks:**

1. Add scene table creation SQL:
```go
const createScenesTable = `
CREATE TABLE IF NOT EXISTS scenes (
    id INTEGER PRIMARY KEY
);
`

const createScenesLEDBarsLEDsTable = `
CREATE TABLE IF NOT EXISTS scenes_ledbars_leds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    ledbar_id INTEGER NOT NULL,
    channel_num INTEGER NOT NULL,
    value INTEGER NOT NULL,
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);
`

const createScenesLEDStripsTable = `
CREATE TABLE IF NOT EXISTS scenes_ledstrips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    red INTEGER NOT NULL,
    green INTEGER NOT NULL,
    blue INTEGER NOT NULL,
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);
`

const createScenesVideoLightsTable = `
CREATE TABLE IF NOT EXISTS scenes_videolights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    on_state INTEGER NOT NULL,
    brightness INTEGER NOT NULL,
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);
`
```

2. Update `InitSchema()` to create scene tables

3. Initialize 4 empty scene slots (IDs 0-3) in scenes table

**Validation:**
- Database file contains new tables
- 4 scene rows exist with IDs 0-3

---

### Phase 2: Storage Interface Updates

**Goal:** Add scene operations to storage interface

**File:** `storage/interface.go`

**Tasks:**

1. Add scene data types:
```go
type SceneData struct {
    LEDStrip    LEDStripState
    LEDBarLEDs  []LEDBarLEDState
    VideoLight1 VideoLightState
    VideoLight2 VideoLightState
}

type LEDStripState struct {
    Red   int
    Green int
    Blue  int
}

type LEDBarLEDState struct {
    LEDBarID   int
    ChannelNum int
    Value      int
}

type VideoLightState struct {
    On         bool
    Brightness int
}
```

2. Add scene methods to Storage interface:
```go
type Storage interface {
    // ... existing methods ...

    // Scene operations
    SceneExists(sceneID int) (bool, error)
    SaveScene(sceneID int, data *SceneData) error
    LoadScene(sceneID int) (*SceneData, error)
    DeleteScene(sceneID int) error
}
```

**Validation:**
- Interface compiles
- Types are defined correctly

---

### Phase 3: Storage Implementation

**Goal:** Implement scene storage operations

**File:** `storage/database.go`

**Tasks:**

1. Implement `SceneExists`:
```go
func (d *Database) SceneExists(sceneID int) (bool, error) {
    // Check if any data exists in scenes_ledstrips for this scene_id
    var count int
    err := d.db.QueryRow(
        "SELECT COUNT(*) FROM scenes_ledstrips WHERE scene_id = ?",
        sceneID,
    ).Scan(&count)
    if err != nil {
        return false, err
    }
    return count > 0, nil
}
```

2. Implement `SaveScene`:
```go
func (d *Database) SaveScene(sceneID int, data *SceneData) error {
    tx, err := d.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Delete existing scene data
    tx.Exec("DELETE FROM scenes_ledbars_leds WHERE scene_id = ?", sceneID)
    tx.Exec("DELETE FROM scenes_ledstrips WHERE scene_id = ?", sceneID)
    tx.Exec("DELETE FROM scenes_videolights WHERE scene_id = ?", sceneID)

    // Insert LED strip state
    _, err = tx.Exec(
        "INSERT INTO scenes_ledstrips (scene_id, red, green, blue) VALUES (?, ?, ?, ?)",
        sceneID, data.LEDStrip.Red, data.LEDStrip.Green, data.LEDStrip.Blue,
    )
    if err != nil {
        return err
    }

    // Insert LED bar LEDs
    for _, led := range data.LEDBarLEDs {
        _, err = tx.Exec(
            "INSERT INTO scenes_ledbars_leds (scene_id, ledbar_id, channel_num, value) VALUES (?, ?, ?, ?)",
            sceneID, led.LEDBarID, led.ChannelNum, led.Value,
        )
        if err != nil {
            return err
        }
    }

    // Insert video light states
    _, err = tx.Exec(
        "INSERT INTO scenes_videolights (scene_id, on_state, brightness) VALUES (?, ?, ?)",
        sceneID, boolToInt(data.VideoLight1.On), data.VideoLight1.Brightness,
    )
    if err != nil {
        return err
    }
    _, err = tx.Exec(
        "INSERT INTO scenes_videolights (scene_id, on_state, brightness) VALUES (?, ?, ?)",
        sceneID, boolToInt(data.VideoLight2.On), data.VideoLight2.Brightness,
    )
    if err != nil {
        return err
    }

    return tx.Commit()
}
```

3. Implement `LoadScene`:
```go
func (d *Database) LoadScene(sceneID int) (*SceneData, error) {
    // Check if scene exists
    exists, err := d.SceneExists(sceneID)
    if err != nil {
        return nil, err
    }
    if !exists {
        return nil, nil  // Empty scene
    }

    data := &SceneData{}

    // Load LED strip
    err = d.db.QueryRow(
        "SELECT red, green, blue FROM scenes_ledstrips WHERE scene_id = ?",
        sceneID,
    ).Scan(&data.LEDStrip.Red, &data.LEDStrip.Green, &data.LEDStrip.Blue)
    if err != nil && err != sql.ErrNoRows {
        return nil, err
    }

    // Load LED bar LEDs
    rows, err := d.db.Query(
        "SELECT ledbar_id, channel_num, value FROM scenes_ledbars_leds WHERE scene_id = ? ORDER BY ledbar_id, channel_num",
        sceneID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var led LEDBarLEDState
        if err := rows.Scan(&led.LEDBarID, &led.ChannelNum, &led.Value); err != nil {
            return nil, err
        }
        data.LEDBarLEDs = append(data.LEDBarLEDs, led)
    }

    // Load video lights
    rows, err = d.db.Query(
        "SELECT on_state, brightness FROM scenes_videolights WHERE scene_id = ? ORDER BY id",
        sceneID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    vlIndex := 0
    for rows.Next() {
        var on int
        var brightness int
        if err := rows.Scan(&on, &brightness); err != nil {
            return nil, err
        }
        if vlIndex == 0 {
            data.VideoLight1 = VideoLightState{On: on != 0, Brightness: brightness}
        } else {
            data.VideoLight2 = VideoLightState{On: on != 0, Brightness: brightness}
        }
        vlIndex++
    }

    return data, nil
}
```

4. Implement `DeleteScene`:
```go
func (d *Database) DeleteScene(sceneID int) error {
    tx, err := d.db.Begin()
    if err != nil {
        return err
    }
    defer tx.Rollback()

    tx.Exec("DELETE FROM scenes_ledbars_leds WHERE scene_id = ?", sceneID)
    tx.Exec("DELETE FROM scenes_ledstrips WHERE scene_id = ?", sceneID)
    tx.Exec("DELETE FROM scenes_videolights WHERE scene_id = ?", sceneID)

    return tx.Commit()
}
```

**Validation:**
- All methods compile
- Unit tests pass for save/load operations

---

### Phase 4: Update Tab Definition

**Goal:** Rename TabFuture2 to TabScenes

**File:** `streamdeck/model.go`

**Tasks:**

1. Update Tab constant:
```go
const (
    TabLightControl Tab = iota
    TabScenes                  // Changed from TabFuture2
    TabFuture3
    TabFuture4
)
```

2. Update Tab.String():
```go
func (t Tab) String() string {
    switch t {
    case TabLightControl:
        return "Lights"
    case TabScenes:
        return "Scenes"
    // ...
    }
}
```

**Validation:**
- Code compiles
- Tab 2 shows "Scenes" label

---

### Phase 5: Add Storage to StreamDeckUI

**Goal:** Give StreamDeckUI access to storage for scene operations

**File:** `streamdeck/model.go`

**Tasks:**

1. Add storage field to StreamDeckUI struct:
```go
type StreamDeckUI struct {
    // ... existing fields ...
    storage storage.Storage
}
```

2. Update NewStreamDeckUI to accept storage:
```go
func NewStreamDeckUI(
    ledStrip *ledstrip.LEDStrip,
    ledBar *ledbar.LEDBar,
    videoLight1 *videolight.VideoLight,
    videoLight2 *videolight.VideoLight,
    storage storage.Storage,  // Add parameter
) (*StreamDeckUI, error) {
    // ...
    ui := &StreamDeckUI{
        // ...
        storage: storage,
    }
    // ...
}
```

3. Update main.go to pass storage to NewStreamDeckUI

**Validation:**
- Code compiles
- StreamDeckUI has access to storage

---

### Phase 6: Implement Scene Operations

**Goal:** Add saveScene and recallScene methods

**File:** `streamdeck/scenes.go` (new file)

**Tasks:**

1. Create new file with scene operations:
```go
package streamdeck

import (
    "log"

    "github.com/kevin/office_lights/storage"
)

// saveScene captures current light state and saves to database
func (s *StreamDeckUI) saveScene(slotIndex int) {
    log.Printf("Saving scene %d...", slotIndex+1)

    // Gather current state from all drivers
    data := &storage.SceneData{
        LEDStrip: storage.LEDStripState{
            Red:   s.ledStrip.R(),
            Green: s.ledStrip.G(),
            Blue:  s.ledStrip.B(),
        },
        VideoLight1: storage.VideoLightState{
            On:         s.videoLight1.IsOn(),
            Brightness: s.videoLight1.Brightness(),
        },
        VideoLight2: storage.VideoLightState{
            On:         s.videoLight2.IsOn(),
            Brightness: s.videoLight2.Brightness(),
        },
    }

    // Gather LED bar state
    data.LEDBarLEDs = s.gatherLEDBarState()

    // Save to database
    if err := s.storage.SaveScene(slotIndex, data); err != nil {
        log.Printf("Error saving scene %d: %v", slotIndex+1, err)
        return
    }

    log.Printf("Scene %d saved successfully", slotIndex+1)

    // Update display
    if err := s.updateTouchscreen(); err != nil {
        log.Printf("Error updating touchscreen: %v", err)
    }
}

// recallScene loads saved state and applies to all lights
func (s *StreamDeckUI) recallScene(slotIndex int) {
    log.Printf("Recalling scene %d...", slotIndex+1)

    // Load from database
    data, err := s.storage.LoadScene(slotIndex)
    if err != nil {
        log.Printf("Error loading scene %d: %v", slotIndex+1, err)
        return
    }

    if data == nil {
        log.Printf("Scene %d is empty", slotIndex+1)
        return
    }

    // Apply to LED strip
    if err := s.ledStrip.SetColor(data.LEDStrip.Red, data.LEDStrip.Green, data.LEDStrip.Blue); err != nil {
        log.Printf("Error setting LED strip: %v", err)
    }

    // Apply to LED bar
    s.applyLEDBarState(data.LEDBarLEDs)

    // Apply to video lights
    if data.VideoLight1.On {
        s.videoLight1.TurnOn(data.VideoLight1.Brightness)
    } else {
        s.videoLight1.TurnOff()
    }
    if data.VideoLight2.On {
        s.videoLight2.TurnOn(data.VideoLight2.Brightness)
    } else {
        s.videoLight2.TurnOff()
    }

    log.Printf("Scene %d recalled successfully", slotIndex+1)
}

// gatherLEDBarState collects all LED bar channel values
func (s *StreamDeckUI) gatherLEDBarState() []storage.LEDBarLEDState {
    var leds []storage.LEDBarLEDState
    // Iterate through all channels and collect values
    // Implementation depends on LEDBar driver API
    return leds
}

// applyLEDBarState applies saved LED bar values
func (s *StreamDeckUI) applyLEDBarState(leds []storage.LEDBarLEDState) {
    // Apply values to LED bar
    // Implementation depends on LEDBar driver API
}
```

**Validation:**
- Methods compile
- Scene save/recall works with test data

---

### Phase 7: Update Event Handlers

**Goal:** Handle button and dial events for Tab 2

**File:** `streamdeck/events.go`

**Tasks:**

1. Update handleButtonPress for TabScenes:
```go
// In handleButtonPress, second row handling:
switch s.currentTab {
case TabLightControl:
    // Existing mode selection
case TabScenes:
    // Recall scene
    s.recallScene(buttonIndex - 4)
default:
    log.Printf("Button %d pressed on unimplemented tab %s", buttonIndex, s.currentTab)
}
```

2. Update handleDialPress for TabScenes:
```go
func (s *StreamDeckUI) handleDialPress(dialIndex int) {
    // ... validation ...

    switch s.currentTab {
    case TabLightControl:
        // Existing toggle behavior
    case TabScenes:
        // Save scene
        s.saveScene(dialIndex)
    default:
        log.Printf("Dial %d pressed on unimplemented tab %s", dialIndex, s.currentTab)
    }
}
```

3. Update handleDialRotate to ignore rotation on TabScenes:
```go
func (s *StreamDeckUI) handleDialRotate(dialIndex int, ticks int) {
    // ...
    if s.currentTab != TabLightControl {
        // Dials only rotate on Tab 1
        return
    }
    // ... existing code ...
}
```

**Validation:**
- Dial click on Tab 2 saves scene
- Button press on Tab 2 recalls scene
- Dial rotation does nothing on Tab 2

---

### Phase 8: Update Rendering

**Goal:** Render Tab 2 UI elements

**File:** `streamdeck/render.go`

**Tasks:**

1. Update renderButton for TabScenes second row:
```go
func (s *StreamDeckUI) renderButton(index int) (image.Image, error) {
    if index < 4 {
        return s.renderTabButton(index)
    }

    switch s.currentTab {
    case TabLightControl:
        return s.renderModeButton(index - 4)
    case TabScenes:
        return s.renderSceneButton(index - 4)
    default:
        return s.renderBlankButton(), nil
    }
}
```

2. Add renderSceneButton:
```go
func (s *StreamDeckUI) renderSceneButton(index int) (image.Image, error) {
    exists, _ := s.storage.SceneExists(index)

    label := fmt.Sprintf("Scene %d", index+1)

    // Different styling for saved vs empty
    return s.renderTextButton(label, exists), nil
}
```

3. Update renderTouchscreen for TabScenes:
```go
func (s *StreamDeckUI) renderTouchscreen() image.Image {
    switch s.currentTab {
    case TabLightControl:
        return s.renderLightControlTouchscreen()
    case TabScenes:
        return s.renderScenesTouchscreen()
    default:
        return s.renderPlaceholderTouchscreen()
    }
}
```

4. Add renderScenesTouchscreen:
```go
func (s *StreamDeckUI) renderScenesTouchscreen() image.Image {
    img := image.NewRGBA(image.Rect(0, 0, touchWidth, touchHeight))
    draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{20, 20, 20, 255}}, image.Point{}, draw.Src)

    for i := 0; i < 4; i++ {
        s.renderSceneSection(img, i)
    }

    return img
}

func (s *StreamDeckUI) renderSceneSection(img *image.RGBA, index int) {
    x := index * sectionWidth
    bounds := image.Rect(x, 0, x+sectionWidth, touchHeight)

    // Background
    bgColor := color.RGBA{40, 40, 40, 255}
    draw.Draw(img, bounds, &image.Uniform{bgColor}, image.Point{x, 0}, draw.Src)

    // Border
    drawVerticalLine(img, x+sectionWidth-1, 0, touchHeight, color.RGBA{80, 80, 80, 255})

    // Label
    label := fmt.Sprintf("Scene %d", index+1)
    drawTextAt(img, label, x+sectionWidth/2, 25, color.RGBA{200, 200, 200, 255}, true)

    // Status
    exists, _ := s.storage.SceneExists(index)
    status := "Empty"
    statusColor := color.RGBA{100, 100, 100, 255}
    if exists {
        status = "Saved"
        statusColor = color.RGBA{100, 200, 100, 255}
    }
    drawTextAt(img, status, x+sectionWidth/2, 60, statusColor, true)
}
```

5. Add "fmt" import if not present

**Validation:**
- Tab 2 shows scene buttons on second row
- Touchscreen shows scene status
- Saved scenes show "Saved" in green
- Empty scenes show "Empty" in gray

---

### Phase 9: Testing

**Goal:** Verify all scene functionality

**Test Cases:**

1. **Database Schema:**
   - [ ] Scene tables created on startup
   - [ ] 4 scene slots exist (IDs 0-3)

2. **Scene Save:**
   - [ ] Dial click on Tab 2 saves scene
   - [ ] All light states captured correctly
   - [ ] Database updated with scene data
   - [ ] Touchscreen shows "Saved" status

3. **Scene Recall:**
   - [ ] Button press on Tab 2 recalls scene
   - [ ] LED strip updated correctly
   - [ ] LED bar updated correctly
   - [ ] Video lights updated correctly
   - [ ] Empty scene shows message, no action

4. **Persistence:**
   - [ ] Saved scene persists after app restart
   - [ ] Scene data loads correctly on recall

5. **UI:**
   - [ ] Tab 2 button shows "Scenes"
   - [ ] Second row shows Scene 1-4 buttons
   - [ ] Saved scenes have different button style
   - [ ] Touchscreen updates after save

6. **Edge Cases:**
   - [ ] Save overwrites existing scene
   - [ ] Recall empty scene does nothing
   - [ ] Multiple rapid saves/recalls work

**Validation:**
- All test cases pass
- No data corruption
- UI updates correctly

---

## Summary of File Changes

| File | Changes |
|------|---------|
| `storage/schema.go` | Add scene table creation SQL |
| `storage/interface.go` | Add SceneData types and scene methods |
| `storage/database.go` | Implement scene storage operations |
| `streamdeck/model.go` | Rename TabFuture2 to TabScenes, add storage field |
| `streamdeck/scenes.go` | New file with saveScene/recallScene |
| `streamdeck/events.go` | Handle Tab 2 button/dial events |
| `streamdeck/render.go` | Render Tab 2 buttons and touchscreen |
| `main.go` | Pass storage to NewStreamDeckUI |

## Implementation Order

1. Phase 1: Database schema (storage/schema.go)
2. Phase 2: Storage interface (storage/interface.go)
3. Phase 3: Storage implementation (storage/database.go)
4. Phase 4: Tab definition update (model.go)
5. Phase 5: Add storage to StreamDeckUI (model.go, main.go)
6. Phase 6: Scene operations (scenes.go)
7. Phase 7: Event handlers (events.go)
8. Phase 8: Rendering (render.go)
9. Phase 9: Testing

## Dependencies

- Storage layer must be complete
- LED bar driver needs method to get/set all channel values
- May need to add `GetAllChannels()` and `SetAllChannels()` to LED bar driver
