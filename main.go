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
	"github.com/kevin/office_lights/storage"
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

	// Get database path from environment variable or use default
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "lights.sqlite3"
	}

	// Create and initialize database
	log.Printf("Opening database at %s...", dbPath)
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	if err := db.InitDefaultData(); err != nil {
		log.Fatalf("Failed to initialize default data: %v", err)
	}

	log.Println("Database initialized successfully")

	// Load state from database
	log.Println("Loading light states from database...")

	// Load LED strip state
	stripR, stripG, stripB, err := db.LoadLEDStripState(0)
	if err != nil {
		log.Printf("Warning: Failed to load LED strip state, using defaults: %v", err)
		stripR, stripG, stripB = 0, 0, 0
	}
	log.Printf("Loaded LED strip state: R=%d, G=%d, B=%d", stripR, stripG, stripB)

	// Load LED bar state
	ledBarChannels, err := db.LoadLEDBarChannels(0)
	if err != nil {
		log.Printf("Warning: Failed to load LED bar state, using defaults: %v", err)
		ledBarChannels = make([]int, 77)
	}
	log.Printf("Loaded LED bar state: %d channels", len(ledBarChannels))

	// Load video light 1 state (database ID 0 -> driver ID 1)
	vl1On, vl1Brightness, err := db.LoadVideoLightState(0)
	if err != nil {
		log.Printf("Warning: Failed to load video light 1 state, using defaults: %v", err)
		vl1On, vl1Brightness = false, 0
	}
	log.Printf("Loaded video light 1 state: on=%v, brightness=%d", vl1On, vl1Brightness)

	// Load video light 2 state (database ID 1 -> driver ID 2)
	vl2On, vl2Brightness, err := db.LoadVideoLightState(1)
	if err != nil {
		log.Printf("Warning: Failed to load video light 2 state, using defaults: %v", err)
		vl2On, vl2Brightness = false, 0
	}
	log.Printf("Loaded video light 2 state: on=%v, brightness=%d", vl2On, vl2Brightness)

	// Instantiate light drivers with loaded state
	log.Println("Initializing light drivers with stored state...")

	// LED Strip
	ledStrip := ledstrip.NewLEDStripWithState(mqttClient, officemqtt.TopicLEDStrip, db, 0, stripR, stripG, stripB)
	log.Println("LED Strip driver initialized")

	// LED Bar
	ledBar, err := ledbar.NewLEDBarWithState(0, mqttClient, officemqtt.TopicLEDBar, db, ledBarChannels)
	if err != nil {
		log.Fatalf("Failed to create LED bar: %v", err)
	}
	log.Println("LED Bar driver initialized")

	// Video Lights
	videoLight1, err := videolight.NewVideoLightWithState(1, mqttClient, officemqtt.TopicVideoLight1, db, vl1On, vl1Brightness)
	if err != nil {
		log.Fatalf("Failed to create video light 1: %v", err)
	}
	log.Println("Video Light 1 driver initialized")

	videoLight2, err := videolight.NewVideoLightWithState(2, mqttClient, officemqtt.TopicVideoLight2, db, vl2On, vl2Brightness)
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

	// Publish initial state to MQTT (sync physical lights with stored state)
	log.Println("Publishing initial state to MQTT...")
	if err := ledStrip.Publish(); err != nil {
		log.Printf("Warning: Failed to publish LED strip initial state: %v", err)
	}
	if err := ledBar.Publish(); err != nil {
		log.Printf("Warning: Failed to publish LED bar initial state: %v", err)
	}
	if err := videoLight1.Publish(); err != nil {
		log.Printf("Warning: Failed to publish video light 1 initial state: %v", err)
	}
	if err := videoLight2.Publish(); err != nil {
		log.Printf("Warning: Failed to publish video light 2 initial state: %v", err)
	}
	log.Println("Initial state published")

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
