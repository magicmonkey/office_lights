# Web Interface Implementation Plan

## Overview

Step-by-step guide for implementing the web-based user interface for the office lights control system.

## Prerequisites

- Completed phases 1-11 (all drivers, MQTT, state storage, TUI)
- Go 1.24.0 or later
- Modern web browser for testing

## Implementation Phases

### Phase 1: Package Structure and Server Setup

**Goal**: Create the web package and basic HTTP server

#### Step 1.1: Create Package Structure

```bash
mkdir -p web/static
```

Create the following files:

```
web/
├── web.go           # Main server setup
├── api.go           # API handlers
├── state.go         # State structures
├── static/
│   ├── index.html   # HTML interface
│   ├── style.css    # Styles
│   └── app.js       # JavaScript
└── web_test.go      # Tests
```

#### Step 1.2: Define State Structures (`state.go`)

```go
package web

import (
	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
)

// RGBW represents a single RGBW LED
type RGBW struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
	W int `json:"w"`
}

// LEDBarSection represents one section of the LED bar
type LEDBarSection struct {
	RGBW  []RGBW `json:"rgbw"`
	White []int  `json:"white"`
}

// LEDBarState represents the complete LED bar state
type LEDBarState struct {
	Section1 LEDBarSection `json:"section1"`
	Section2 LEDBarSection `json:"section2"`
}

// LEDStripState represents the LED strip state
type LEDStripState struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

// VideoLightState represents a video light state
type VideoLightState struct {
	On         bool `json:"on"`
	Brightness int  `json:"brightness"`
}

// State represents the complete system state
type State struct {
	LEDStrip    LEDStripState   `json:"ledStrip"`
	LEDBar      LEDBarState     `json:"ledBar"`
	VideoLight1 VideoLightState `json:"videoLight1"`
	VideoLight2 VideoLightState `json:"videoLight2"`
}

// BuildState reads current state from all drivers
func BuildState(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) (*State, error) {
	state := &State{}

	// LED Strip
	r, g, b := strip.GetColor()
	state.LEDStrip = LEDStripState{R: r, G: g, B: b}

	// LED Bar - need to read all LEDs in both sections
	state.LEDBar.Section1.RGBW = make([]RGBW, 6)
	state.LEDBar.Section1.White = make([]int, 13)
	state.LEDBar.Section2.RGBW = make([]RGBW, 6)
	state.LEDBar.Section2.White = make([]int, 13)

	// Section 1 RGBW
	for i := 0; i < 6; i++ {
		r, g, b, w, err := bar.GetRGBW(1, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section1.RGBW[i] = RGBW{R: r, G: g, B: b, W: w}
	}

	// Section 1 White
	for i := 0; i < 13; i++ {
		val, err := bar.GetWhite(1, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section1.White[i] = val
	}

	// Section 2 RGBW
	for i := 0; i < 6; i++ {
		r, g, b, w, err := bar.GetRGBW(2, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section2.RGBW[i] = RGBW{R: r, G: g, B: b, W: w}
	}

	// Section 2 White
	for i := 0; i < 13; i++ {
		val, err := bar.GetWhite(2, i)
		if err != nil {
			return nil, err
		}
		state.LEDBar.Section2.White[i] = val
	}

	// Video Lights
	on1, brightness1 := vl1.GetState()
	state.VideoLight1 = VideoLightState{On: on1, Brightness: brightness1}

	on2, brightness2 := vl2.GetState()
	state.VideoLight2 = VideoLightState{On: on2, Brightness: brightness2}

	return state, nil
}

// ApplyState applies state to all drivers
func ApplyState(
	state *State,
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) error {
	// LED Strip
	if err := strip.SetColor(state.LEDStrip.R, state.LEDStrip.G, state.LEDStrip.B); err != nil {
		return err
	}

	// LED Bar Section 1 RGBW
	for i, rgbw := range state.LEDBar.Section1.RGBW {
		if err := bar.SetRGBW(1, i, rgbw.R, rgbw.G, rgbw.B, rgbw.W); err != nil {
			return err
		}
	}

	// LED Bar Section 1 White
	for i, val := range state.LEDBar.Section1.White {
		if err := bar.SetWhite(1, i, val); err != nil {
			return err
		}
	}

	// LED Bar Section 2 RGBW
	for i, rgbw := range state.LEDBar.Section2.RGBW {
		if err := bar.SetRGBW(2, i, rgbw.R, rgbw.G, rgbw.B, rgbw.W); err != nil {
			return err
		}
	}

	// LED Bar Section 2 White
	for i, val := range state.LEDBar.Section2.White {
		if err := bar.SetWhite(2, i, val); err != nil {
			return err
		}
	}

	// Video Light 1
	if state.VideoLight1.On {
		if err := vl1.TurnOn(state.VideoLight1.Brightness); err != nil {
			return err
		}
	} else {
		if err := vl1.TurnOff(); err != nil {
			return err
		}
	}

	// Video Light 2
	if state.VideoLight2.On {
		if err := vl2.TurnOn(state.VideoLight2.Brightness); err != nil {
			return err
		}
	} else {
		if err := vl2.TurnOff(); err != nil {
			return err
		}
	}

	return nil
}

// Validate checks if state values are within valid ranges
func (s *State) Validate() error {
	// LED Strip
	if s.LEDStrip.R < 0 || s.LEDStrip.R > 255 {
		return fmt.Errorf("LED strip R value out of range: %d", s.LEDStrip.R)
	}
	if s.LEDStrip.G < 0 || s.LEDStrip.G > 255 {
		return fmt.Errorf("LED strip G value out of range: %d", s.LEDStrip.G)
	}
	if s.LEDStrip.B < 0 || s.LEDStrip.B > 255 {
		return fmt.Errorf("LED strip B value out of range: %d", s.LEDStrip.B)
	}

	// LED Bar - validate RGBW values
	validateRGBW := func(section string, rgbwList []RGBW) error {
		for i, rgbw := range rgbwList {
			if rgbw.R < 0 || rgbw.R > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] R out of range: %d", section, i, rgbw.R)
			}
			if rgbw.G < 0 || rgbw.G > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] G out of range: %d", section, i, rgbw.G)
			}
			if rgbw.B < 0 || rgbw.B > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] B out of range: %d", section, i, rgbw.B)
			}
			if rgbw.W < 0 || rgbw.W > 255 {
				return fmt.Errorf("LED bar %s RGBW[%d] W out of range: %d", section, i, rgbw.W)
			}
		}
		return nil
	}

	if err := validateRGBW("section1", s.LEDBar.Section1.RGBW); err != nil {
		return err
	}
	if err := validateRGBW("section2", s.LEDBar.Section2.RGBW); err != nil {
		return err
	}

	// LED Bar - validate white values
	validateWhite := func(section string, white []int) error {
		for i, val := range white {
			if val < 0 || val > 255 {
				return fmt.Errorf("LED bar %s white[%d] out of range: %d", section, i, val)
			}
		}
		return nil
	}

	if err := validateWhite("section1", s.LEDBar.Section1.White); err != nil {
		return err
	}
	if err := validateWhite("section2", s.LEDBar.Section2.White); err != nil {
		return err
	}

	// Video Lights
	if s.VideoLight1.Brightness < 0 || s.VideoLight1.Brightness > 100 {
		return fmt.Errorf("video light 1 brightness out of range: %d", s.VideoLight1.Brightness)
	}
	if s.VideoLight2.Brightness < 0 || s.VideoLight2.Brightness > 100 {
		return fmt.Errorf("video light 2 brightness out of range: %d", s.VideoLight2.Brightness)
	}

	return nil
}
```

#### Step 1.3: Create Server Structure (`web.go`)

```go
package web

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/kevin/office_lights/drivers/ledbar"
	"github.com/kevin/office_lights/drivers/ledstrip"
	"github.com/kevin/office_lights/drivers/videolight"
)

//go:embed static/*
var staticFiles embed.FS

// Server represents the web server
type Server struct {
	ledStrip    *ledstrip.LEDStrip
	ledBar      *ledbar.LEDBar
	videoLight1 *videolight.VideoLight
	videoLight2 *videolight.VideoLight
	httpServer  *http.Server
	mu          sync.Mutex // Protect concurrent access
}

// NewServer creates a new web server
func NewServer(
	strip *ledstrip.LEDStrip,
	bar *ledbar.LEDBar,
	vl1 *videolight.VideoLight,
	vl2 *videolight.VideoLight,
) *Server {
	return &Server{
		ledStrip:    strip,
		ledBar:      bar,
		videoLight1: vl1,
		videoLight2: vl2,
	}
}

// Start starts the HTTP server
func (s *Server) Start(port string) error {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/", http.FileServer(http.FS(staticFiles)))

	// API endpoints
	mux.HandleFunc("/api", s.handleAPI)
	mux.HandleFunc("/health", s.handleHealth)

	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Web server starting on http://localhost:%s", port)
	return s.httpServer.ListenAndServe()
}

// Stop gracefully stops the HTTP server
func (s *Server) Stop() error {
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	return nil
}
```

#### Verification

- Package structure created
- State structures defined
- Server skeleton implemented
- Code compiles

---

### Phase 2: API Implementation

**Goal**: Implement GET and POST API endpoints

#### Step 2.1: Implement API Handler (`api.go`)

```go
package web

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// handleAPI handles both GET and POST requests to /api
func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Enable CORS for development
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		s.handleGetState(w, r)
	case "POST":
		s.handlePostState(w, r)
	default:
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

// handleGetState returns the current state as JSON
func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	state, err := BuildState(s.ledStrip, s.ledBar, s.videoLight1, s.videoLight2)
	if err != nil {
		log.Printf("Error building state: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to read state: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(state); err != nil {
		log.Printf("Error encoding state: %v", err)
		http.Error(w, `{"error":"Failed to encode state"}`, http.StatusInternalServerError)
		return
	}
}

// handlePostState applies the provided state
func (s *Server) handlePostState(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, `{"error":"Failed to read request body"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON
	var state State
	if err := json.Unmarshal(body, &state); err != nil {
		log.Printf("Error parsing JSON: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Invalid JSON: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Validate state
	if err := state.Validate(); err != nil {
		log.Printf("Validation error: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Validation failed: %v"}`, err), http.StatusBadRequest)
		return
	}

	// Apply state to drivers
	if err := ApplyState(&state, s.ledStrip, s.ledBar, s.videoLight1, s.videoLight2); err != nil {
		log.Printf("Error applying state: %v", err)
		http.Error(w, fmt.Sprintf(`{"error":"Failed to apply state: %v"}`, err), http.StatusInternalServerError)
		return
	}

	// Return the updated state
	updatedState, err := BuildState(s.ledStrip, s.ledBar, s.videoLight1, s.videoLight2)
	if err != nil {
		log.Printf("Error building updated state: %v", err)
		http.Error(w, `{"error":"State applied but failed to read back"}`, http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(updatedState); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		return
	}

	log.Println("State updated successfully")
}

// handleHealth returns a simple health check response
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
```

#### Verification

- GET /api returns JSON state
- POST /api accepts and applies state
- Validation works correctly
- Errors are handled properly
- CORS headers set for development

---

### Phase 3: HTML Interface

**Goal**: Create the HTML structure and basic layout

#### Step 3.1: Create HTML (`static/index.html`)

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Office Lights Control</title>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body>
    <div class="container">
        <header>
            <h1>Office Lights Control</h1>
            <div class="status">
                <span id="connection-status">Connecting...</span>
                <span id="last-update">Never</span>
            </div>
        </header>

        <main class="grid">
            <!-- LED Strip -->
            <section class="card">
                <h2>LED Strip</h2>
                <div class="control-group">
                    <label for="strip-r">Red (0-255)</label>
                    <input type="range" id="strip-r" min="0" max="255" value="0">
                    <span id="strip-r-value">0</span>
                </div>
                <div class="control-group">
                    <label for="strip-g">Green (0-255)</label>
                    <input type="range" id="strip-g" min="0" max="255" value="0">
                    <span id="strip-g-value">0</span>
                </div>
                <div class="control-group">
                    <label for="strip-b">Blue (0-255)</label>
                    <input type="range" id="strip-b" min="0" max="255" value="0">
                    <span id="strip-b-value">0</span>
                </div>
                <div class="control-group">
                    <label for="strip-color">Color Picker</label>
                    <input type="color" id="strip-color" value="#000000">
                </div>
                <div class="preview" id="strip-preview"></div>
            </section>

            <!-- LED Bar -->
            <section class="card">
                <h2>LED Bar</h2>
                <div class="control-group">
                    <label>Section</label>
                    <div class="button-group">
                        <button id="ledbar-section-1" class="active">Section 1</button>
                        <button id="ledbar-section-2">Section 2</button>
                    </div>
                </div>
                <div class="control-group">
                    <label>Mode</label>
                    <div class="button-group">
                        <button id="ledbar-mode-rgbw" class="active">RGBW</button>
                        <button id="ledbar-mode-white">White</button>
                    </div>
                </div>
                <div id="ledbar-rgbw-controls">
                    <div class="control-group">
                        <label for="ledbar-led">LED (1-6)</label>
                        <input type="number" id="ledbar-led" min="1" max="6" value="1">
                    </div>
                    <div class="control-group">
                        <label for="ledbar-r">Red (0-255)</label>
                        <input type="range" id="ledbar-r" min="0" max="255" value="0">
                        <span id="ledbar-r-value">0</span>
                    </div>
                    <div class="control-group">
                        <label for="ledbar-g">Green (0-255)</label>
                        <input type="range" id="ledbar-g" min="0" max="255" value="0">
                        <span id="ledbar-g-value">0</span>
                    </div>
                    <div class="control-group">
                        <label for="ledbar-b">Blue (0-255)</label>
                        <input type="range" id="ledbar-b" min="0" max="255" value="0">
                        <span id="ledbar-b-value">0</span>
                    </div>
                    <div class="control-group">
                        <label for="ledbar-w">White (0-255)</label>
                        <input type="range" id="ledbar-w" min="0" max="255" value="0">
                        <span id="ledbar-w-value">0</span>
                    </div>
                </div>
                <div id="ledbar-white-controls" style="display: none;">
                    <div class="control-group">
                        <label for="ledbar-white-led">LED (1-13)</label>
                        <input type="number" id="ledbar-white-led" min="1" max="13" value="1">
                    </div>
                    <div class="control-group">
                        <label for="ledbar-white">Brightness (0-255)</label>
                        <input type="range" id="ledbar-white" min="0" max="255" value="0">
                        <span id="ledbar-white-value">0</span>
                    </div>
                </div>
            </section>

            <!-- Video Light 1 -->
            <section class="card">
                <h2>Video Light 1</h2>
                <div class="control-group">
                    <label for="vl1-on">Power</label>
                    <label class="switch">
                        <input type="checkbox" id="vl1-on">
                        <span class="slider"></span>
                    </label>
                </div>
                <div class="control-group">
                    <label for="vl1-brightness">Brightness (0-100)</label>
                    <input type="range" id="vl1-brightness" min="0" max="100" value="0">
                    <span id="vl1-brightness-value">0</span>
                </div>
                <div class="indicator" id="vl1-indicator"></div>
            </section>

            <!-- Video Light 2 -->
            <section class="card">
                <h2>Video Light 2</h2>
                <div class="control-group">
                    <label for="vl2-on">Power</label>
                    <label class="switch">
                        <input type="checkbox" id="vl2-on">
                        <span class="slider"></span>
                    </label>
                </div>
                <div class="control-group">
                    <label for="vl2-brightness">Brightness (0-100)</label>
                    <input type="range" id="vl2-brightness" min="0" max="100" value="0">
                    <span id="vl2-brightness-value">0</span>
                </div>
                <div class="indicator" id="vl2-indicator"></div>
            </section>
        </main>

        <footer>
            <div id="error-message" class="error" style="display: none;"></div>
        </footer>
    </div>

    <script src="/static/app.js"></script>
</body>
</html>
```

#### Verification

- HTML loads in browser
- All controls render correctly
- Layout is responsive

---

### Phase 4: CSS Styling

**Goal**: Style the interface for usability and aesthetics

#### Step 4.1: Create Stylesheet (`static/style.css`)

```css
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    background: #1a1a1a;
    color: #ffffff;
    line-height: 1.6;
}

.container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 20px;
}

header {
    text-align: center;
    margin-bottom: 30px;
    padding-bottom: 20px;
    border-bottom: 2px solid #333;
}

h1 {
    color: #00ADD8;
    margin-bottom: 10px;
}

.status {
    display: flex;
    justify-content: center;
    gap: 20px;
    font-size: 0.9em;
    color: #888;
}

#connection-status.connected {
    color: #00FF00;
}

#connection-status.disconnected {
    color: #FF0000;
}

.grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 20px;
    margin-bottom: 20px;
}

.card {
    background: #2a2a2a;
    border-radius: 10px;
    padding: 20px;
    border: 2px solid #333;
    transition: border-color 0.3s;
}

.card:hover {
    border-color: #00ADD8;
}

.card h2 {
    color: #00ADD8;
    margin-bottom: 15px;
    font-size: 1.2em;
}

.control-group {
    margin-bottom: 15px;
}

.control-group label {
    display: block;
    margin-bottom: 5px;
    color: #ccc;
    font-size: 0.9em;
}

input[type="range"] {
    width: 100%;
    height: 6px;
    background: #444;
    border-radius: 3px;
    outline: none;
    -webkit-appearance: none;
}

input[type="range"]::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    background: #00ADD8;
    border-radius: 50%;
    cursor: pointer;
}

input[type="range"]::-moz-range-thumb {
    width: 18px;
    height: 18px;
    background: #00ADD8;
    border-radius: 50%;
    cursor: pointer;
    border: none;
}

input[type="number"] {
    width: 100%;
    padding: 8px;
    background: #333;
    border: 1px solid #555;
    border-radius: 5px;
    color: #fff;
    font-size: 1em;
}

input[type="color"] {
    width: 100%;
    height: 40px;
    border: none;
    border-radius: 5px;
    cursor: pointer;
}

.button-group {
    display: flex;
    gap: 10px;
}

.button-group button {
    flex: 1;
    padding: 8px 15px;
    background: #333;
    border: 2px solid #555;
    border-radius: 5px;
    color: #fff;
    cursor: pointer;
    transition: all 0.3s;
}

.button-group button.active {
    background: #00ADD8;
    border-color: #00ADD8;
    color: #000;
}

.button-group button:hover {
    border-color: #00ADD8;
}

/* Toggle Switch */
.switch {
    position: relative;
    display: inline-block;
    width: 60px;
    height: 34px;
}

.switch input {
    opacity: 0;
    width: 0;
    height: 0;
}

.slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: #555;
    transition: 0.4s;
    border-radius: 34px;
}

.slider:before {
    position: absolute;
    content: "";
    height: 26px;
    width: 26px;
    left: 4px;
    bottom: 4px;
    background-color: white;
    transition: 0.4s;
    border-radius: 50%;
}

input:checked + .slider {
    background-color: #00ADD8;
}

input:checked + .slider:before {
    transform: translateX(26px);
}

/* Preview boxes */
.preview, .indicator {
    width: 100%;
    height: 50px;
    border-radius: 5px;
    margin-top: 15px;
    border: 2px solid #444;
    transition: background-color 0.3s;
}

.indicator {
    background: #333;
}

.indicator.on {
    box-shadow: 0 0 20px rgba(0, 173, 216, 0.6);
}

/* Value display */
span[id$="-value"] {
    display: inline-block;
    min-width: 40px;
    text-align: right;
    color: #FDDD00;
    font-weight: bold;
}

/* Error message */
.error {
    background: #ff4444;
    color: white;
    padding: 15px;
    border-radius: 5px;
    margin-top: 20px;
    text-align: center;
}

footer {
    margin-top: 20px;
}

/* Responsive */
@media (max-width: 768px) {
    .grid {
        grid-template-columns: 1fr;
    }

    .container {
        padding: 10px;
    }
}
```

#### Verification

- Interface looks clean and modern
- Controls are easy to use
- Responsive on mobile devices

---

### Phase 5: JavaScript Implementation

**Goal**: Implement the frontend logic

#### Step 5.1: Create JavaScript (`static/app.js`)

(Due to length, will provide key sections)

```javascript
// State management
let currentState = null;
let updateTimer = null;
let pollingInterval = null;
let isUpdating = false;

// LED Bar UI state
let selectedSection = 1;
let selectedMode = 'rgbw'; // 'rgbw' or 'white'

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    initializeEventListeners();
    startPolling();
    loadState();
});

// Load current state from API
async function loadState() {
    try {
        const response = await fetch('/api');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const state = await response.json();
        currentState = state;
        updateUI(state);
        updateConnectionStatus(true);
        updateLastUpdateTime();
    } catch (error) {
        console.error('Error loading state:', error);
        updateConnectionStatus(false);
        showError('Failed to load state: ' + error.message);
    }
}

// Save state to API
async function saveState(state) {
    if (isUpdating) return;

    isUpdating = true;
    try {
        const response = await fetch('/api', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(state),
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || 'Unknown error');
        }

        const updatedState = await response.json();
        currentState = updatedState;
        updateConnectionStatus(true);
        hideError();
    } catch (error) {
        console.error('Error saving state:', error);
        updateConnectionStatus(false);
        showError('Failed to save state: ' + error.message);
    } finally {
        isUpdating = false;
    }
}

// Debounced update
function scheduleUpdate() {
    if (updateTimer) {
        clearTimeout(updateTimer);
    }
    updateTimer = setTimeout(() => {
        if (currentState) {
            saveState(currentState);
        }
    }, 300); // 300ms debounce
}

// Polling
function startPolling() {
    pollingInterval = setInterval(loadState, 3000); // Poll every 3 seconds
}

function stopPolling() {
    if (pollingInterval) {
        clearInterval(pollingInterval);
    }
}

// Update UI from state
function updateUI(state) {
    // LED Strip
    document.getElementById('strip-r').value = state.ledStrip.r;
    document.getElementById('strip-g').value = state.ledStrip.g;
    document.getElementById('strip-b').value = state.ledStrip.b;
    updateStripValueDisplays();
    updateStripColorPicker();
    updateStripPreview();

    // LED Bar - update based on current selection
    updateLEDBarUI(state);

    // Video Light 1
    document.getElementById('vl1-on').checked = state.videoLight1.on;
    document.getElementById('vl1-brightness').value = state.videoLight1.brightness;
    updateVideoLight1ValueDisplay();
    updateVideoLight1Indicator();

    // Video Light 2
    document.getElementById('vl2-on').checked = state.videoLight2.on;
    document.getElementById('vl2-brightness').value = state.videoLight2.brightness;
    updateVideoLight2ValueDisplay();
    updateVideoLight2Indicator();
}

// Initialize event listeners
function initializeEventListeners() {
    // LED Strip
    ['strip-r', 'strip-g', 'strip-b'].forEach(id => {
        document.getElementById(id).addEventListener('input', (e) => {
            if (!currentState) return;
            const channel = id.split('-')[1];
            currentState.ledStrip[channel] = parseInt(e.target.value);
            updateStripValueDisplays();
            updateStripColorPicker();
            updateStripPreview();
            scheduleUpdate();
        });
    });

    document.getElementById('strip-color').addEventListener('input', (e) => {
        if (!currentState) return;
        const rgb = hexToRgb(e.target.value);
        currentState.ledStrip.r = rgb.r;
        currentState.ledStrip.g = rgb.g;
        currentState.ledStrip.b = rgb.b;
        updateUI(currentState);
        scheduleUpdate();
    });

    // LED Bar section buttons
    document.getElementById('ledbar-section-1').addEventListener('click', () => {
        selectedSection = 1;
        document.getElementById('ledbar-section-1').classList.add('active');
        document.getElementById('ledbar-section-2').classList.remove('active');
        if (currentState) updateLEDBarUI(currentState);
    });

    document.getElementById('ledbar-section-2').addEventListener('click', () => {
        selectedSection = 2;
        document.getElementById('ledbar-section-2').classList.add('active');
        document.getElementById('ledbar-section-1').classList.remove('active');
        if (currentState) updateLEDBarUI(currentState);
    });

    // LED Bar mode buttons
    document.getElementById('ledbar-mode-rgbw').addEventListener('click', () => {
        selectedMode = 'rgbw';
        document.getElementById('ledbar-mode-rgbw').classList.add('active');
        document.getElementById('ledbar-mode-white').classList.remove('active');
        document.getElementById('ledbar-rgbw-controls').style.display = 'block';
        document.getElementById('ledbar-white-controls').style.display = 'none';
        if (currentState) updateLEDBarUI(currentState);
    });

    document.getElementById('ledbar-mode-white').addEventListener('click', () => {
        selectedMode = 'white';
        document.getElementById('ledbar-mode-white').classList.add('active');
        document.getElementById('ledbar-mode-rgbw').classList.remove('active');
        document.getElementById('ledbar-rgbw-controls').style.display = 'none';
        document.getElementById('ledbar-white-controls').style.display = 'block';
        if (currentState) updateLEDBarUI(currentState);
    });

    // LED Bar RGBW controls
    // ... (similar to LED strip)

    // Video Light 1
    document.getElementById('vl1-on').addEventListener('change', (e) => {
        if (!currentState) return;
        currentState.videoLight1.on = e.target.checked;
        updateVideoLight1Indicator();
        scheduleUpdate();
    });

    document.getElementById('vl1-brightness').addEventListener('input', (e) => {
        if (!currentState) return;
        currentState.videoLight1.brightness = parseInt(e.target.value);
        updateVideoLight1ValueDisplay();
        scheduleUpdate();
    });

    // Video Light 2
    // ... (similar to Video Light 1)
}

// Helper functions
function hexToRgb(hex) {
    const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
        r: parseInt(result[1], 16),
        g: parseInt(result[2], 16),
        b: parseInt(result[3], 16)
    } : null;
}

function rgbToHex(r, g, b) {
    return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

function updateConnectionStatus(connected) {
    const status = document.getElementById('connection-status');
    if (connected) {
        status.textContent = 'Connected';
        status.className = 'connected';
    } else {
        status.textContent = 'Disconnected';
        status.className = 'disconnected';
    }
}

function updateLastUpdateTime() {
    document.getElementById('last-update').textContent = 'Just now';
}

function showError(message) {
    const errorDiv = document.getElementById('error-message');
    errorDiv.textContent = message;
    errorDiv.style.display = 'block';
}

function hideError() {
    document.getElementById('error-message').style.display = 'none';
}

// Additional helper functions for UI updates
// ...
```

---

### Phase 6: Main Integration

**Goal**: Integrate web server with main application

#### Step 6.1: Update main.go

Add after TUI section:

```go
// Check if web mode is requested
useWeb := false
webPort := os.Getenv("WEB_PORT")
if webPort == "" {
    webPort = "8080"
}

if os.Getenv("WEB") != "" {
    useWeb = true
}
for _, arg := range os.Args[1:] {
    if arg == "web" {
        useWeb = true
    }
}

if useWeb {
    log.Printf("Starting web server on port %s...", webPort)
    webServer := web.NewServer(ledStrip, ledBar, videoLight1, videoLight2)
    go func() {
        if err := webServer.Start(webPort); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Web server error: %v", err)
        }
    }()
    log.Printf("Web interface available at http://localhost:%s", webPort)

    // If only web mode, don't exit - wait for signal
    if !useTUI {
        // Web mode only - wait for shutdown signal
        // (signal handling code already exists below)
    }
}
```

#### Verification

- Can run with `./office_lights web`
- Can run with `WEB=1 ./office_lights`
- Web server starts on port 8080
- Can access http://localhost:8080

---

### Phase 7: Testing

**Goal**: Test all functionality

#### Unit Tests

Create `web/web_test.go`:

```go
package web

import (
    "testing"
)

func TestStateValidation(t *testing.T) {
    // Test valid state
    // Test invalid RGB values
    // Test invalid brightness values
}

func TestBuildState(t *testing.T) {
    // Test with mock drivers
}

func TestApplyState(t *testing.T) {
    // Test with mock drivers
}
```

#### Integration Tests

- Start server with test drivers
- Make GET request
- Make POST request
- Verify state changes

#### Manual Testing Checklist

- [ ] Web interface loads
- [ ] LED Strip controls work
- [ ] LED Bar RGBW controls work
- [ ] LED Bar White controls work
- [ ] Video Light 1 controls work
- [ ] Video Light 2 controls work
- [ ] Color picker updates RGB sliders
- [ ] RGB sliders update color picker
- [ ] Toggle switches work
- [ ] State persists across page reload
- [ ] Multiple browsers can control
- [ ] Error messages display correctly
- [ ] Mobile layout works
- [ ] Changes publish to MQTT
- [ ] Changes save to database

---

## Implementation Order Summary

1. **Phase 1**: Package structure and server setup
2. **Phase 2**: API implementation
3. **Phase 3**: HTML structure
4. **Phase 4**: CSS styling
5. **Phase 5**: JavaScript logic
6. **Phase 6**: Main integration
7. **Phase 7**: Testing

## Success Criteria

- Web server starts and serves interface
- API returns correct JSON state
- API accepts and applies state updates
- HTML interface is functional and responsive
- All controls work as expected
- Changes are published to MQTT
- Changes are saved to database
- Multiple browsers can control simultaneously
- Error handling works correctly

## Estimated Complexity

- **Simple**: Phases 1-2 (backend structure)
- **Moderate**: Phases 3-4 (HTML/CSS)
- **Moderate-Complex**: Phase 5 (JavaScript)
- **Simple**: Phases 6-7 (integration and testing)

## Notes

- Use embedded files (`embed.FS`) for static assets
- Mutex protects concurrent API access
- Debouncing prevents excessive MQTT messages
- Polling keeps UI synchronized
- No external JavaScript dependencies
- Works offline (after initial load)
