package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
	officemqtt "github.com/kevin/office_lights/mqtt"
)

// Lights holds all light driver instances
type Lights struct {
	LEDStrip    *ledstrip.LEDStrip
	LEDBar      *ledbar.LEDBar
	VideoLight1 *videolight.VideoLight
	VideoLight2 *videolight.VideoLight
}

func main() {
	log.Println("Office Lights Control System Starting...")

	// Get MQTT broker address from environment variable or use default
	broker := os.Getenv("MQTT_BROKER")
	if broker == "" {
		broker = "tcp://localhost:1883"
	}

	// Get MQTT client ID from environment variable or use default
	clientID := os.Getenv("MQTT_CLIENT_ID")
	if clientID == "" {
		clientID = "office_lights_controller"
	}

	// Create MQTT client configuration
	config := officemqtt.Config{
		Broker:   broker,
		ClientID: clientID,
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
	}

	// Create and connect MQTT client
	mqttClient, err := officemqtt.NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create MQTT client: %v", err)
	}

	if err := mqttClient.Connect(); err != nil {
		log.Fatalf("Failed to connect to MQTT broker: %v", err)
	}
	defer mqttClient.Disconnect()

	log.Println("MQTT client connected successfully")

	// Instantiate light drivers
	log.Println("Initializing light drivers...")

	// LED Strip
	ledStrip := ledstrip.NewLEDStrip(mqttClient, officemqtt.TopicLEDStrip)
	log.Println("LED Strip driver initialized")

	// LED Bar
	ledBar, err := ledbar.NewLEDBar(0, mqttClient, officemqtt.TopicLEDBar)
	if err != nil {
		log.Fatalf("Failed to create LED bar: %v", err)
	}
	log.Println("LED Bar driver initialized")

	// Video Lights
	videoLight1, err := videolight.NewVideoLight(1, mqttClient, officemqtt.TopicVideoLight1)
	if err != nil {
		log.Fatalf("Failed to create video light 1: %v", err)
	}
	log.Println("Video Light 1 driver initialized")

	videoLight2, err := videolight.NewVideoLight(2, mqttClient, officemqtt.TopicVideoLight2)
	if err != nil {
		log.Fatalf("Failed to create video light 2: %v", err)
	}
	log.Println("Video Light 2 driver initialized")

	// Store all lights for easy access
	lights := &Lights{
		LEDStrip:    ledStrip,
		LEDBar:      ledBar,
		VideoLight1: videoLight1,
		VideoLight2: videoLight2,
	}

	// Demonstrate basic functionality (optional - can be removed)
	demonstrateLights(lights)

	log.Println("Office Lights Control System Ready")

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down gracefully...", sig)

	// Turn off all lights before shutdown
	log.Println("Turning off all lights...")
	if err := lights.LEDStrip.TurnOff(); err != nil {
		log.Printf("Error turning off LED strip: %v", err)
	}
	if err := lights.LEDBar.TurnOffAll(); err != nil {
		log.Printf("Error turning off LED bar: %v", err)
	}
	if err := lights.VideoLight1.TurnOff(); err != nil {
		log.Printf("Error turning off video light 1: %v", err)
	}
	if err := lights.VideoLight2.TurnOff(); err != nil {
		log.Printf("Error turning off video light 2: %v", err)
	}

	// Cleanup will happen via defer statements
	log.Println("Shutdown complete")
}

// demonstrateLights shows basic functionality of all light types
// This is optional and can be removed or disabled via environment variable
func demonstrateLights(lights *Lights) {
	// Skip demo if environment variable is set
	if os.Getenv("SKIP_DEMO") != "" {
		return
	}

	log.Println("Running light demonstration...")

	// Demo LED Strip - set to a warm white
	log.Println("Demo: Setting LED strip to warm white")
	if err := lights.LEDStrip.SetColor(255, 200, 150); err != nil {
		log.Printf("Error setting LED strip color: %v", err)
	}

	// Demo LED Bar - set first RGBW LED in section 1 to blue
	log.Println("Demo: Setting LED bar first RGBW to blue")
	if err := lights.LEDBar.SetRGBW(1, 0, 0, 0, 255, 100); err != nil {
		log.Printf("Error setting LED bar RGBW: %v", err)
	}

	// Demo Video Lights - turn on at 75% brightness
	log.Println("Demo: Turning on video lights at 75% brightness")
	if err := lights.VideoLight1.TurnOn(75); err != nil {
		log.Printf("Error turning on video light 1: %v", err)
	}
	if err := lights.VideoLight2.TurnOn(75); err != nil {
		log.Printf("Error turning on video light 2: %v", err)
	}

	log.Println("Light demonstration complete")
}
