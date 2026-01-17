# Web Interface Architecture

## Overview

Implement a web-based user interface that provides browser-based control for all office lights. The interface runs as an HTTP server in a separate goroutine and provides both an API endpoint and an interactive HTML interface.

## Requirements

### Architecture

The web interface consists of two main components:

1. **REST API** (`/api`)
   - GET: Returns complete light status as JSON
   - POST: Accepts complete light status as JSON and applies changes

2. **HTML Interface** (`/`)
   - Single-page application
   - AJAX communication with the API
   - Real-time UI updates
   - Interactive controls for all lights

### HTTP Server

**Technology:** Standard library `net/http`
- Simple, no external dependencies for basic functionality
- Runs in a dedicated goroutine
- Does not block main application

**Port:** Configurable via environment variable (default: 8080)

**Endpoints:**
- `GET /` - Serve HTML interface
- `GET /api` - Get current light status (JSON)
- `POST /api` - Update light status (JSON)
- `GET /health` - Health check endpoint (optional)

### API Design

#### GET /api Response

Returns the complete state of all lights:

```json
{
  "ledStrip": {
    "r": 255,
    "g": 200,
    "b": 150
  },
  "ledBar": {
    "section1": {
      "rgbw": [
        {"r": 0, "g": 0, "b": 255, "w": 100},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0}
      ],
      "white": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    },
    "section2": {
      "rgbw": [
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0},
        {"r": 0, "g": 0, "b": 0, "w": 0}
      ],
      "white": [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
    }
  },
  "videoLight1": {
    "on": true,
    "brightness": 75
  },
  "videoLight2": {
    "on": false,
    "brightness": 0
  }
}
```

#### POST /api Request

Accepts the same JSON structure as GET response. Updates all lights to the specified state.

**Behavior:**
- Validates JSON structure
- Validates value ranges (RGB 0-255, brightness 0-100, etc.)
- Applies changes to drivers (which publish to MQTT and save to DB)
- Returns updated state or error

**Response:**
- Success (200): Same JSON structure as GET
- Error (400): `{"error": "description of error"}`
- Error (500): `{"error": "internal server error"}`

### HTML Interface Design

#### Technology Stack

**Vanilla JavaScript** - No framework required for simplicity
- Fetch API for AJAX requests
- DOM manipulation for UI updates
- Event listeners for user input

**CSS** - Simple, responsive design
- Flexbox/Grid layout
- Color pickers for RGB controls
- Sliders for numeric values
- Toggle switches for on/off

#### Layout

```
┌─────────────────────────────────────────────────┐
│  Office Lights Control                          │
├─────────────────────────────────────────────────┤
│                                                 │
│  ┌──────────────┐  ┌──────────────────────────┐│
│  │ LED Strip    │  │ LED Bar                  ││
│  │ R: [slider]  │  │ Section: [1] [2]         ││
│  │ G: [slider]  │  │ Mode: [RGBW] [White]     ││
│  │ B: [slider]  │  │ LED: [1-6] or [1-13]     ││
│  │ [color pick] │  │ R: [slider]              ││
│  └──────────────┘  │ G: [slider]              ││
│                    │ B: [slider]              ││
│  ┌──────────────┐  │ W: [slider]              ││
│  │ Video Light 1│  └──────────────────────────┘│
│  │ [ON/OFF]     │                              │
│  │ Bright: [slider]│  ┌──────────────┐         │
│  └──────────────┘  │ Video Light 2│         │
│                    │ [ON/OFF]     │         │
│                    │ Bright: [slider]│         │
│                    └──────────────┘         │
│                                                 │
│  Status: Connected | Last update: 2s ago       │
└─────────────────────────────────────────────────┘
```

#### User Interactions

**LED Strip:**
- RGB sliders (0-255 each)
- Color picker (HTML5 `<input type="color">`)
- Preset color buttons (optional)

**LED Bar:**
- Section selector (1 or 2)
- Mode selector (RGBW or White)
- LED index selector
- RGBW sliders when in RGBW mode
- Single brightness slider when in White mode
- Visual representation of all LEDs

**Video Lights:**
- Toggle switch for on/off
- Brightness slider (0-100)
- Visual indicator (color changes when on)

#### Update Strategy

**Polling:**
- Poll GET /api every 2-5 seconds to stay synchronized
- Update UI only if values changed (avoid flickering)
- Show "last updated" timestamp

**User Changes:**
- Debounce rapid changes (wait 300ms after last change)
- POST complete state to /api
- Update UI with response to confirm
- Show loading/saving indicator

**Error Handling:**
- Display error messages in status bar
- Retry on network failure (with exponential backoff)
- Highlight problematic controls
- Maintain last known good state

### Package Structure

```
web/
├── web.go           # HTTP server setup and handlers
├── api.go           # API endpoint handlers
├── handlers.go      # HTTP handler functions
├── state.go         # State structure and marshaling
├── static/
│   ├── index.html   # Main HTML page
│   ├── style.css    # Styles
│   └── app.js       # JavaScript application
└── web_test.go      # Unit tests
```

### Integration with Existing Code

#### Driver Access

The web server needs access to all driver instances:

```go
type Server struct {
    ledStrip    *ledstrip.LEDStrip
    ledBar      *ledbar.LEDBar
    videoLight1 *videolight.VideoLight
    videoLight2 *videolight.VideoLight
    httpServer  *http.Server
}
```

**Reading State:**
- Use existing getter methods: `GetColor()`, `GetState()`, `GetRGBW()`, `GetWhite()`
- Build JSON response from current driver state

**Updating State:**
- Parse incoming JSON
- Validate all values
- Call driver methods: `SetColor()`, `SetRGBW()`, `SetWhite()`, `TurnOn()`, `TurnOff()`
- Drivers handle MQTT publishing and database persistence

#### Main Integration

```go
func main() {
    // ... existing setup ...

    // Check if web mode is requested
    useWeb := os.Getenv("WEB") != "" || hasArg("web")
    webPort := os.Getenv("WEB_PORT")
    if webPort == "" {
        webPort = "8080"
    }

    if useWeb {
        webServer := web.NewServer(ledStrip, ledBar, videoLight1, videoLight2)
        go webServer.Start(webPort)
        log.Printf("Web interface started on http://localhost:%s", webPort)
    }

    // ... rest of main ...
}
```

### Concurrency Considerations

**Goroutine Safety:**
- HTTP server runs in separate goroutine
- Multiple requests may access drivers concurrently
- Need mutex protection for driver state access

**Options:**
1. **Add mutex to drivers** - Protect internal state
2. **Add mutex to web server** - Serialize API requests
3. **Use channels** - Send updates via channels

**Recommended:** Add mutex to web server layer (option 2)
- Simpler to implement
- Drivers remain unchanged
- API requests naturally serialized

### Error Handling

**Network Errors:**
- Connection failures
- Timeout handling
- Graceful degradation

**Validation Errors:**
- Invalid JSON structure
- Out-of-range values
- Missing required fields

**Driver Errors:**
- MQTT publish failures
- Database save failures
- Driver-specific errors

**Response Format:**
```json
{
  "error": "error description",
  "field": "ledStrip.r",  // optional: which field caused error
  "details": "additional context"  // optional
}
```

### Security Considerations

**Basic Security:**
- No authentication in initial version (local network use)
- CORS headers (if needed for development)
- Input validation and sanitization
- Rate limiting (optional)

**Future Enhancements:**
- Basic auth (username/password)
- HTTPS support
- API tokens
- CSRF protection

### Testing Strategy

#### Unit Tests

**API Tests:**
- Test JSON marshaling/unmarshaling
- Test state building from drivers
- Test state application to drivers
- Test validation logic

**Handler Tests:**
- Test GET /api returns correct JSON
- Test POST /api with valid data
- Test POST /api with invalid data
- Test error responses

#### Integration Tests

- Test with mock drivers
- Test concurrent requests
- Test state synchronization
- Test error propagation

#### Manual Testing

- Browser compatibility (Chrome, Firefox, Safari)
- Mobile responsiveness
- Network failure scenarios
- Rapid control changes
- Multiple browser windows

### Performance Considerations

**Efficiency:**
- Minimize state reads (cache if needed)
- Debounce rapid changes
- Use conditional updates
- Gzip compression for responses

**Scalability:**
- Single driver instance (not a concern)
- Handle multiple concurrent browsers
- Memory footprint minimal

### Monitoring and Logging

**Logging:**
- HTTP requests (method, path, status)
- API errors
- Driver update errors
- Startup/shutdown events

**Metrics (optional):**
- Request count
- Error rate
- Response times
- Active connections

### Deployment

**Standalone Mode:**
```bash
./office_lights web
# or
WEB=1 ./office_lights
```

**With Custom Port:**
```bash
WEB_PORT=3000 ./office_lights web
```

**With TUI (Not Recommended):**
- Running both TUI and Web simultaneously not recommended
- TUI would still suppress logs
- Web runs in background, TUI in foreground

### Future Enhancements

**Not in Initial Implementation:**
- WebSocket support for real-time updates
- Scene presets (save/load configurations)
- Scheduling (time-based control)
- Multiple LED bar support
- Animations and effects
- Mobile app (native or PWA)
- Voice control integration

## Technology Decisions

### Why Standard Library HTTP Server?

**Pros:**
- No external dependencies
- Simple and well-documented
- Sufficient for our needs
- Good performance

**Cons:**
- No built-in routing (but we only need 2-3 routes)
- Manual CORS handling
- No middleware system

**Decision:** Use standard library. Simple needs don't justify a framework.

### Why Vanilla JavaScript?

**Pros:**
- No build step required
- Simple deployment (static files)
- Fast loading
- Easy to understand

**Cons:**
- More verbose than frameworks
- Manual DOM manipulation
- No reactive data binding

**Decision:** Vanilla JS is sufficient. Interface is simple enough that a framework would be overkill.

### Why Polling Instead of WebSockets?

**Pros:**
- Simpler implementation
- Works with any HTTP server
- Easier to debug
- Lower complexity

**Cons:**
- Not real-time (2-5 second delay)
- More bandwidth (but negligible for our use)

**Decision:** Start with polling. Can add WebSockets later if needed.

## Success Criteria

- Web server starts successfully
- HTML interface loads in browser
- GET /api returns current state
- POST /api updates lights correctly
- Changes reflected in physical lights (via MQTT)
- Changes persisted to database
- Multiple browsers can control simultaneously
- Error messages displayed clearly
- Works on desktop and mobile browsers
