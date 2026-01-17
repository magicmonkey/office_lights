# Implementation Order

## Overview
Recommended order for implementing the office lights control system to ensure smooth development and testing.

## Phase 1: Foundation
1. **Project Setup** (spec/01-project-setup.md)
   - Initialize Go module
   - Set up project structure
   - Install dependencies

2. **MQTT Infrastructure** (spec/02-mqtt-infrastructure.md)
   - Set up MQTT client
   - Test basic connectivity
   - Implement publish functionality

## Phase 2: Individual Drivers
Implement each driver independently. They can be developed in parallel or in any order:

3. **LED Strip Driver** (spec/03-ledstrip-driver.md)
   - Simplest driver - good starting point
   - Test JSON message formatting

4. **Video Light Driver** (spec/05-videolight-driver.md)
   - Simple message format
   - Multiple instances needed

5. **LED Bar Driver** (spec/04-ledbar-driver.md)
   - Most complex message format
   - Requires careful testing of CSV generation

## Phase 3: Integration
6. **Main Orchestration** (spec/06-main-orchestration.md)
   - Bring all drivers together
   - Create light instances
   - Test end-to-end functionality

## Phase 4: Quality Assurance
7. **Testing** (spec/07-testing-strategy.md)
   - Write unit tests for each driver
   - Integration tests
   - Manual testing with actual lights

## Phase 5: Future Enhancements
8. **UI Integration** (spec/08-future-ui-integration.md)
   - Plan and implement user interface
   - Create API layer
   - Add presets and scenes

## Development Tips

### For Each Driver
1. Create the package structure
2. Define the state struct
3. Implement constructor
4. Implement state management methods
5. Implement message formatting
6. Implement publish method
7. Write tests
8. Test with actual hardware (if available)

### Testing Strategy
- Test message formatting before testing with actual MQTT
- Use a local MQTT broker (like Mosquitto) for development
- Subscribe to topics with an MQTT client to verify messages
- Start with mock/stub implementations before connecting to real lights

### Validation Checklist
Before moving to the next phase, ensure:
- [ ] Code compiles without errors
- [ ] All tests pass
- [ ] Message formats are verified against specification
- [ ] Error handling is in place
- [ ] Code is documented

## Estimated Complexity
- **Low Complexity:** Project setup, MQTT infrastructure, LED strip, video light
- **Medium Complexity:** Main orchestration, testing
- **High Complexity:** LED bar (due to complex message format), UI integration

## Dependencies
- Each driver depends on MQTT infrastructure
- Main orchestration depends on all drivers
- Testing can be done incrementally as each driver is completed
- UI integration depends on main orchestration being complete
