# Configuration

## Environment Variables

The office lights control system can be configured using the following environment variables:

### MQTT Broker Connection

- `MQTT_BROKER` - MQTT broker address (default: `tcp://localhost:1883`)
  - Examples:
    - `tcp://localhost:1883`
    - `tcp://192.168.1.100:1883`
    - `ssl://broker.example.com:8883`

- `MQTT_CLIENT_ID` - Unique client identifier (default: `office_lights_controller`)
  - Should be unique if running multiple instances

- `MQTT_USERNAME` - MQTT broker username (optional)
  - Only needed if your broker requires authentication

- `MQTT_PASSWORD` - MQTT broker password (optional)
  - Only needed if your broker requires authentication

### State Storage

- `DB_PATH` - Path to SQLite database file (default: `lights.sqlite3`)
  - Stores persistent state for all lights
  - Will be created automatically on first run
  - Example: `./data/lights.db`

## Example Usage

### Basic (Local MQTT Broker)
```bash
./office_lights
```

### With Custom Broker
```bash
export MQTT_BROKER="tcp://192.168.1.100:1883"
./office_lights
```

### With Authentication
```bash
export MQTT_BROKER="tcp://broker.example.com:1883"
export MQTT_USERNAME="myuser"
export MQTT_PASSWORD="mypassword"
./office_lights
```

### All Options
```bash
export MQTT_BROKER="tcp://192.168.1.100:1883"
export MQTT_CLIENT_ID="office_lights_main"
export MQTT_USERNAME="admin"
export MQTT_PASSWORD="secret"
export DB_PATH="./lights.sqlite3"
./office_lights
```

## Building

```bash
go build
```

## Running

```bash
./office_lights
```

To stop the application, press `Ctrl+C` for graceful shutdown.

## MQTT Topics

The following topics are used:

- `kevinoffice/ledstrip/sequence` - LED strip control
- `kevinoffice/ledbar/0` - LED bar control
- `kevinoffice/videolight/1/command/light:0` - Video light 1 control
- `kevinoffice/videolight/2/command/light:0` - Video light 2 control

## Testing MQTT Connection

To test the MQTT connection, you can use a tool like `mosquitto_sub` to subscribe to topics:

```bash
mosquitto_sub -h localhost -t 'kevinoffice/#' -v
```

This will show all messages published to topics under `kevinoffice/`.

## State Persistence

### Database File

The application stores the current state of all lights in a SQLite database file. This allows lights to restore their previous state after the application restarts.

**Database Location:** `lights.sqlite3` (current directory by default)

**What's stored:**
- LED strip RGB values
- LED bar channel values (all 77 channels)
- Video light on/off state and brightness

### Database Backup

To backup your light states:
```bash
# Simple backup (when application is stopped)
cp lights.sqlite3 lights.sqlite3.backup

# Or while running
sqlite3 lights.sqlite3 ".backup lights.sqlite3.backup"
```

### Restore from Backup

```bash
cp lights.sqlite3.backup lights.sqlite3
./office_lights
```

### Reset to Defaults

To reset all lights to their default state (off):
```bash
# Stop the application
rm lights.sqlite3
# Restart - a new database will be created with defaults
./office_lights
```

### Database Inspection

You can inspect the database contents using the SQLite CLI:
```bash
sqlite3 lights.sqlite3

# List tables
.tables

# View LED strip state
SELECT * FROM ledstrips;

# View video lights state
SELECT * FROM videolights;

# View LED bar channels
SELECT * FROM ledbars_leds ORDER BY channel_num;

# Exit
.quit
```
