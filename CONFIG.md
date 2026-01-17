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

### Application Behavior

- `SKIP_DEMO` - Skip the light demonstration on startup (optional)
  - When not set, the application will demonstrate each light type on startup
  - Set to any value to skip the demo

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
./office_lights
```

### Skip Demonstration Mode
```bash
export SKIP_DEMO=1
./office_lights
```

By default, the application runs a brief demonstration of each light type on startup. This helps verify connectivity and shows you what each light does. To skip this demonstration, set the `SKIP_DEMO` environment variable.

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
