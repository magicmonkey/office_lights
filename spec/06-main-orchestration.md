# Main Orchestration

## Overview
Implement the main.go file that ties everything together and manages all light instances.

## Tasks

### 1. Import Required Packages
- Import all driver packages
- Import MQTT client library
- Import standard libraries (fmt, log, os, etc.)

### 2. Initialize MQTT Client
- Set up MQTT broker configuration
  - Read from environment variables or config file
  - Broker address, port, client ID
- Create and connect MQTT client
- Handle connection errors

### 3. Instantiate Light Drivers
Create instances for:
- 1 LED strip: `NewLEDStrip(mqttClient)`
- 1 LED bar: `NewLEDBar(0, mqttClient)` (based on topic "kevinoffice/ledbar/0")
- 2 Video lights:
  - `NewVideoLight(1, mqttClient)`
  - `NewVideoLight(2, mqttClient)`

### 4. Create Light Manager/Registry
- Store all light instances in a manageable structure
- Provide methods to access specific lights
- Consider using a map or struct to organize lights by type and ID

### 5. Implement Graceful Shutdown
- Handle OS signals (SIGINT, SIGTERM)
- Disconnect MQTT client cleanly
- Optional: Turn off all lights on shutdown

### 6. Add Basic Demonstration/Testing Code
- For initial validation, add code to:
  - Set LED strip to a color
  - Set LED bar values
  - Turn on video lights
- This can be removed when UI is implemented

### 7. Configuration Management
- Consider using environment variables for:
  - MQTT broker address
  - MQTT broker port
  - Client ID
- Or use a config file (JSON, YAML, TOML)

### 8. Logging
- Set up logging for:
  - Application startup
  - MQTT connection status
  - Light state changes
  - Errors

### 9. Future UI Integration Preparation
- Design interface that UI will use to interact with lights
- Consider creating a simple API or command handler
- Document how future UI will trigger light changes

## Success Criteria
- Application starts successfully
- MQTT client connects to broker
- All light driver instances are created
- Can send commands to all light types
- Application shuts down cleanly
- Proper error handling and logging throughout
