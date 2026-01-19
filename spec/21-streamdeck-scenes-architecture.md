# Stream Deck+ Scenes Architecture (Tab 2)

## Overview

This specification defines the architecture for Tab 2 of the Stream Deck+ interface, which provides functionality to save and recall lighting "scenes" (presets). A scene captures the complete state of all lights and can be stored in one of 4 slots.

## Feature Summary

- **4 Scene Slots**: Each slot can store a complete snapshot of all light states
- **Save Scene**: Click a dial to save the current light state to that slot
- **Recall Scene**: Press a second-row button to apply the saved scene to the lights
- **Persistent Storage**: Scenes are stored in the SQLite database

## Database Schema

### New Tables

```sql
-- Scene metadata (4 rows, one per slot)
CREATE TABLE scenes (
    id INTEGER PRIMARY KEY
);

-- Scene LED bar state (mirrors ledbars_leds structure)
CREATE TABLE scenes_ledbars_leds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    ledbar_id INTEGER NOT NULL,
    channel_num INTEGER NOT NULL,
    value INTEGER NOT NULL,
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);

-- Scene LED strip state (mirrors ledstrips structure)
CREATE TABLE scenes_ledstrips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    red INTEGER NOT NULL,
    green INTEGER NOT NULL,
    blue INTEGER NOT NULL,
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);

-- Scene video light state (mirrors videolights structure)
CREATE TABLE scenes_videolights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    on_state INTEGER NOT NULL,
    brightness INTEGER NOT NULL,
    FOREIGN KEY (scene_id) REFERENCES scenes(id)
);
```

### Initial Data

On first run, create 4 empty scene slots (IDs 0-3). Empty scenes have no associated data in the child tables.

### Data Relationships

```
scenes (id: 0-3)
    └── scenes_ledbars_leds (scene_id -> scenes.id)
    └── scenes_ledstrips (scene_id -> scenes.id)
    └── scenes_videolights (scene_id -> scenes.id)
```

## User Interface

### Tab 2 Layout

```
Top Row (Tab Selection):
┌──────────┬──────────┬──────────┬──────────┐
│  Tab 1   │  Tab 2   │  Tab 3   │  Tab 4   │
│ (Lights) │ (Scenes) │ (Future) │ (Future) │
└──────────┴──────────┴──────────┴──────────┘

Second Row (Scene Slots):
┌──────────┬──────────┬──────────┬──────────┐
│ Scene 1  │ Scene 2  │ Scene 3  │ Scene 4  │
│          │          │          │          │
└──────────┴──────────┴──────────┴──────────┘

Touchscreen:
┌─────────────────────────────────────────────┐
│  Scene 1  │  Scene 2  │  Scene 3  │  Scene 4│
│  [status] │  [status] │  [status] │ [status]│
└─────────────────────────────────────────────┘

Dials:
   Dial 1      Dial 2      Dial 3      Dial 4
   (Save 1)    (Save 2)    (Save 3)    (Save 4)
```

### Button Behavior (Second Row)

| Button | Action |
|--------|--------|
| Button 4 | Recall Scene 1 |
| Button 5 | Recall Scene 2 |
| Button 6 | Recall Scene 3 |
| Button 7 | Recall Scene 4 |

### Dial Behavior

| Dial | Rotation | Click |
|------|----------|-------|
| Dial 1 | No effect | Save current state to Scene 1 |
| Dial 2 | No effect | Save current state to Scene 2 |
| Dial 3 | No effect | Save current state to Scene 3 |
| Dial 4 | No effect | Save current state to Scene 4 |

### Touchscreen Display

Each of the 4 sections shows:
- Scene label ("Scene 1", "Scene 2", etc.)
- Status indicator:
  - "Empty" if no scene is saved
  - "Saved" if a scene exists
  - Brief feedback on save/recall ("Saving...", "Loading...")

### Visual Feedback

**On Save (dial click):**
1. Show "Saving..." on touchscreen section
2. Save to database
3. Show "Saved" confirmation
4. Return to normal display

**On Recall (button press):**
1. Show "Loading..." on touchscreen section
2. Load from database
3. Apply to all lights
4. Return to normal display (or show "Empty" if no scene)

## Technical Architecture

### State Model Updates

```go
// Add to model.go or create scenes.go

// SceneSlot represents a saved scene
type SceneSlot struct {
    ID       int
    HasData  bool  // Whether scene has been saved
}

// SceneData holds the complete state for a scene
type SceneData struct {
    LEDStrip    LEDStripState
    LEDBar      LEDBarState
    VideoLight1 VideoLightState
    VideoLight2 VideoLightState
}

type LEDStripState struct {
    Red   int
    Green int
    Blue  int
}

type LEDBarState struct {
    Channels []int  // All channel values
}

type VideoLightState struct {
    On         bool
    Brightness int
}
```

### Storage Layer Updates

Add to storage interface:

```go
// Scene storage operations
type SceneStorage interface {
    // Check if a scene slot has data
    SceneExists(sceneID int) (bool, error)

    // Save current light state to a scene slot
    SaveScene(sceneID int, data SceneData) error

    // Load scene data from a slot
    LoadScene(sceneID int) (*SceneData, error)

    // Delete a scene (optional, for future use)
    DeleteScene(sceneID int) error
}
```

### Event Handling

```go
// In events.go, update dial press handler for Tab 2
func (s *StreamDeckUI) handleDialPress(dialIndex int) {
    switch s.currentTab {
    case TabLightControl:
        // Existing toggle behavior
    case TabScenes:
        s.saveScene(dialIndex)
    }
}

// In events.go, update button press handler for Tab 2
func (s *StreamDeckUI) handleButtonPress(buttonIndex int) {
    // ... tab selection for top row ...

    // Second row handling
    switch s.currentTab {
    case TabLightControl:
        // Mode selection
    case TabScenes:
        s.recallScene(buttonIndex - 4)
    }
}
```

### Scene Operations

```go
// saveScene captures current light state and saves to database
func (s *StreamDeckUI) saveScene(slotIndex int) {
    // 1. Gather current state from all drivers
    data := SceneData{
        LEDStrip: LEDStripState{
            Red:   s.ledStrip.R(),
            Green: s.ledStrip.G(),
            Blue:  s.ledStrip.B(),
        },
        LEDBar: LEDBarState{
            Channels: s.ledBar.GetAllChannels(),
        },
        VideoLight1: VideoLightState{
            On:         s.videoLight1.IsOn(),
            Brightness: s.videoLight1.Brightness(),
        },
        VideoLight2: VideoLightState{
            On:         s.videoLight2.IsOn(),
            Brightness: s.videoLight2.Brightness(),
        },
    }

    // 2. Save to database
    if err := s.storage.SaveScene(slotIndex, data); err != nil {
        log.Printf("Error saving scene %d: %v", slotIndex, err)
        return
    }

    // 3. Update display
    log.Printf("Scene %d saved", slotIndex+1)
}

// recallScene loads saved state and applies to all lights
func (s *StreamDeckUI) recallScene(slotIndex int) {
    // 1. Load from database
    data, err := s.storage.LoadScene(slotIndex)
    if err != nil {
        log.Printf("Error loading scene %d: %v", slotIndex, err)
        return
    }

    if data == nil {
        log.Printf("Scene %d is empty", slotIndex+1)
        return
    }

    // 2. Apply to LED strip
    s.ledStrip.SetColor(data.LEDStrip.Red, data.LEDStrip.Green, data.LEDStrip.Blue)

    // 3. Apply to LED bar
    s.ledBar.SetAllChannels(data.LEDBar.Channels)

    // 4. Apply to video lights
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

    log.Printf("Scene %d recalled", slotIndex+1)
}
```

### Rendering

```go
// renderScenesTouchscreen renders Tab 2 touchscreen
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

    // Draw section background
    // ...

    // Draw label
    label := fmt.Sprintf("Scene %d", index+1)
    drawTextAt(img, label, x+sectionWidth/2, 20, color.White, true)

    // Draw status
    exists, _ := s.storage.SceneExists(index)
    status := "Empty"
    if exists {
        status = "Saved"
    }
    drawTextAt(img, status, x+sectionWidth/2, 60, color.Gray, true)
}

// renderSceneButton renders a scene slot button
func (s *StreamDeckUI) renderSceneButton(index int) (image.Image, error) {
    exists, _ := s.storage.SceneExists(index)

    label := fmt.Sprintf("Scene %d", index+1)
    if exists {
        // Show as available/saved
        return s.renderTextButton(label, true), nil
    } else {
        // Show as empty
        return s.renderTextButton(label, false), nil
    }
}
```

## Dependencies

### Storage Layer

The StreamDeckUI needs access to the storage layer for scene persistence:

```go
type StreamDeckUI struct {
    // ... existing fields ...
    storage Storage  // Add storage interface reference
}
```

### Driver Methods

May need to add methods to drivers for bulk state operations:

- `ledBar.GetAllChannels() []int` - Get all channel values
- `ledBar.SetAllChannels([]int)` - Set all channel values at once

## Error Handling

- **Database errors**: Log error, show brief error message on touchscreen
- **Empty scene recall**: Show "Empty" status, no action taken
- **Driver errors**: Log error, continue with other lights

## Testing

### Unit Tests

- Scene save captures correct state
- Scene load returns correct data
- Empty scene handling
- Database operations work correctly

### Integration Tests

- Save scene, restart app, recall scene
- Save/recall with various light configurations
- Multiple scene slots work independently

### Manual Tests

- [ ] Tab 2 button shows "Scenes" label
- [ ] Second row shows Scene 1-4 buttons
- [ ] Touchscreen shows scene status
- [ ] Dial click saves current light state
- [ ] Button press recalls saved scene
- [ ] Empty scenes show "Empty" status
- [ ] Saved scenes show "Saved" status
- [ ] Scene persists after app restart

## Future Enhancements

1. **Scene names**: Allow custom names for scenes
2. **Scene preview**: Show thumbnail/summary of scene contents
3. **Scene copy**: Copy scene from one slot to another
4. **Scene delete**: Clear a scene slot
5. **Scene export/import**: Backup and restore scenes
