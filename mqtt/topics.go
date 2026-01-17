package mqtt

// MQTT topic constants for office lights
const (
	// TopicLEDStrip is the topic for controlling the RGB LED strip
	TopicLEDStrip = "kevinoffice/ledstrip/sequence"

	// TopicLEDBar is the topic for controlling LED bar 0
	TopicLEDBar = "kevinoffice/ledbar/0"

	// TopicVideoLight1 is the topic for controlling video light 1
	TopicVideoLight1 = "kevinoffice/videolight/1/command/light:0"

	// TopicVideoLight2 is the topic for controlling video light 2
	TopicVideoLight2 = "kevinoffice/videolight/2/command/light:0"
)
