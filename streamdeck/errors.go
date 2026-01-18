package streamdeck

import "errors"

var (
	// ErrNoDevice is returned when no Stream Deck device is found
	ErrNoDevice = errors.New("no Stream Deck device found")

	// ErrInvalidButton is returned when an invalid button index is used
	ErrInvalidButton = errors.New("invalid button index")

	// ErrInvalidDial is returned when an invalid dial index is used
	ErrInvalidDial = errors.New("invalid dial index")

	// ErrInvalidSection is returned when an invalid section index is used
	ErrInvalidSection = errors.New("invalid section index")
)
