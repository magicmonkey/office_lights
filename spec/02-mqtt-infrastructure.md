# MQTT Infrastructure

## Overview
Set up the MQTT client connection and message publishing infrastructure.

## Tasks

### 1. MQTT Client Configuration
- Define MQTT broker connection parameters (host, port, client ID)
- Create connection options structure
- Implement connection retry logic
- Handle connection loss/reconnection

### 2. MQTT Client Wrapper
- Create a simple wrapper or manager for MQTT operations
- Implement `Connect()` method
- Implement `Disconnect()` method
- Implement `Publish(topic string, payload interface{})` method

### 3. Topic Management
- Define constants or configuration for MQTT topics:
  - `kevinoffice/ledstrip/sequence`
  - `kevinoffice/ledbar/0`
  - `kevinoffice/videolight/1/command/light:0`
  - `kevinoffice/videolight/2/command/light:0`

### 4. Error Handling
- Handle publish failures
- Log MQTT connection events
- Implement timeout mechanisms

### 5. Testing Considerations
- Consider how to mock MQTT for testing
- Test connection establishment
- Test message publishing

## Success Criteria
- Can successfully connect to an MQTT broker
- Can publish messages to topics
- Connection handles disconnects gracefully
- Error conditions are properly logged
