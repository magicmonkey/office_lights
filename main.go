package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
	officemqtt "github.com/kevin/office_lights/mqtt"
	"github.com/kevin/office_lights/storage"
	"github.com/kevin/office_lights/tui"
	"github.com/kevin/office_lights/web"
)

func main() {
	// Check if TUI mode is requested early to suppress logs
	useTUI := false
	if len(os.Args) > 1 && os.Args[1] == "tui" {
		useTUI = true
	}
	if os.Getenv("TUI") != "" {
		useTUI = true
	}

	// Check if web mode is requested
	useWeb := false
	if len(os.Args) > 1 && os.Args[1] == "web" {
		useWeb = true
	}
	if os.Getenv("WEB") != "" {
		useWeb = true
	}

	// Disable logging if TUI mode is active (web mode still shows logs)
	if useTUI {
		log.SetOutput(io.Discard)
	}

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

	// Check if database already has data
	hasData, err := db.HasData()
	if err != nil {
		log.Fatalf("Failed to check for existing data: %v", err)
	}

	if hasData {
		log.Println("Database already contains data, skipping initialization")
	} else {
		log.Println("Database is empty, initializing with default data")
		if err := db.InitDefaultData(); err != nil {
			log.Fatalf("Failed to initialize default data: %v", err)
		}
	}

	log.Println("Database ready")

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

	log.Println("Office Lights Control System Ready")

	// Start TUI in a goroutine if requested
	if useTUI {
		go func() {
			log.Println("Starting TUI mode...")
			if err := tui.Run(ledStrip, ledBar, videoLight1, videoLight2); err != nil {
				log.Fatalf("TUI error: %v", err)
			}
			log.Println("TUI exited")
		}()
	}

	// Start web server in a goroutine if requested
	if useWeb {
		// Get web server port from environment variable or use default
		port := os.Getenv("WEB_PORT")
		if port == "" {
			port = "8080"
		}

		// Create and start web server
		webServer := web.NewServer(ledStrip, ledBar, videoLight1, videoLight2)

		// Start web server in a goroutine so it doesn't block
		go func() {
			log.Printf("Starting web interface on port %s...", port)
			if err := webServer.Start(port); err != nil {
				log.Fatalf("Web server error: %v", err)
			}
		}()

		log.Printf("Web interface available at http://localhost:%s", port)
	}

	// If no UI is requested, just note that we're running in headless mode
	if !useTUI && !useWeb {
		log.Println("Running in headless mode (no UI). Press Ctrl+C to exit.")
	}

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for termination signal
	sig := <-sigChan
	log.Printf("Received signal %v, shutting down gracefully...", sig)

	// Cleanup will happen via defer statements
	log.Println("Shutdown complete")
}
