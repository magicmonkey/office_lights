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
