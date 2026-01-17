# Project Setup and Structure

## Overview
Initial project setup for the office lights control system in Go.

## Tasks

### 1. Initialize Go Module
- Create or verify `go.mod` file with appropriate module name
- Set Go version (recommend Go 1.21 or later)

### 2. Define Project Structure
```
office_lights/
├── main.go
├── drivers/
│   ├── ledstrip/
│   ├── ledbar/
│   └── videolight/
├── spec/
└── README.md
```

### 3. Install Dependencies
- MQTT client library (e.g., `github.com/eclipse/paho.mqtt.golang`)
- JSON encoding/decoding (standard library `encoding/json`)
- Any additional utilities needed

### 4. Create Basic main.go
- Package declaration
- Import statements
- Placeholder main function
- Basic error handling structure

## Success Criteria
- `go mod tidy` runs without errors
- Project structure is in place
- Can run `go build` successfully (even with minimal code)
