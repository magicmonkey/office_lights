package storage

// StateStore defines the interface for persistent state storage
type StateStore interface {
	// SaveLEDStripState saves the RGB state for an LED strip
	SaveLEDStripState(id int, r, g, b int) error

	// LoadLEDStripState loads the RGB state for an LED strip
	LoadLEDStripState(id int) (r, g, b int, err error)

	// SaveLEDBarChannels saves all 77 channel values for an LED bar
	SaveLEDBarChannels(ledbarID int, channels []int) error

	// LoadLEDBarChannels loads all 77 channel values for an LED bar
	LoadLEDBarChannels(ledbarID int) ([]int, error)

	// SaveVideoLightState saves the on/off and brightness state for a video light
	SaveVideoLightState(id int, on bool, brightness int) error

	// LoadVideoLightState loads the on/off and brightness state for a video light
	LoadVideoLightState(id int) (on bool, brightness int, err error)
}

// SceneStore defines the interface for scene storage operations
type SceneStore interface {
	// SceneExists checks if a scene slot has saved data
	SceneExists(sceneID int) (bool, error)

	// GetSceneName returns the name of a scene slot
	GetSceneName(sceneID int) (string, error)

	// GetSceneBgColor returns the background color of a scene slot (hex string like "#FF5500")
	GetSceneBgColor(sceneID int) (string, error)

	// SaveScene saves the current light state to a scene slot
	SaveScene(sceneID int, data *SceneData) error

	// LoadScene loads scene data from a slot (returns nil if empty)
	LoadScene(sceneID int) (*SceneData, error)

	// DeleteScene clears a scene slot
	DeleteScene(sceneID int) error
}

// SceneData holds the complete state for a saved scene
type SceneData struct {
	LEDStrip    LEDStripState
	LEDBarLEDs  []LEDBarLEDState
	VideoLights []VideoLightState
}

// LEDStripState holds the RGB state of an LED strip
type LEDStripState struct {
	Red   int
	Green int
	Blue  int
}

// LEDBarLEDState holds a single LED bar channel value
type LEDBarLEDState struct {
	LEDBarID   int
	ChannelNum int
	Value      int
}

// VideoLightState holds the state of a video light
type VideoLightState struct {
	ID         int
	On         bool
	Brightness int
}
