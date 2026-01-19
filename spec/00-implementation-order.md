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

## Phase 5: State Persistence
8. **State Storage Infrastructure** (spec/09-state-storage.md)
   - SQLite database setup
   - Schema creation
   - Storage layer implementation
   - Load and save operations

9. **Driver State Integration** (spec/10-driver-state-integration.md)
   - Update drivers to use storage layer
   - Add state persistence to all drivers
   - Maintain backward compatibility

See **spec/11-state-storage-implementation-order.md** for detailed implementation steps.

## Phase 6: Text User Interface
10. **TUI Architecture** (spec/13-tui-architecture.md)
    - Design terminal-based user interface
    - Choose TUI framework (Bubbletea)
    - Define component structure
    - Plan keyboard navigation

11. **TUI Implementation** (spec/14-tui-implementation-plan.md)
    - Install dependencies (Bubbletea, Bubbles, Lipgloss)
    - Implement component models (LED strip, LED bar, video lights)
    - Create view and update functions
    - Add keyboard handling
    - Integrate with main application

## Phase 7: Web User Interface
12. **Web Interface Architecture** (spec/15-web-interface-architecture.md)
    - Design REST API endpoints
    - Define JSON state structure
    - Plan HTML/CSS/JavaScript interface
    - Choose technology stack (vanilla JS, standard library HTTP)
    - Design concurrency handling

13. **Web Interface Implementation** (spec/16-web-interface-implementation-plan.md)
    - Create web package with HTTP server
    - Implement GET/POST API endpoints
    - Build HTML interface with controls
    - Add CSS styling
    - Implement JavaScript logic with polling
    - Integrate with main application

## Phase 8: Stream Deck+ Interface
14. **Stream Deck Architecture** (spec/17-streamdeck-architecture.md)
    - Design Stream Deck+ hardware interface
    - Define 4 operational modes
    - Plan button, touchscreen, and dial controls
    - Choose USB HID library (github.com/muesli/streamdeck)
    - Design image rendering pipeline
    - Plan event handling architecture

15. **Stream Deck Implementation** (spec/18-streamdeck-implementation-plan.md)
    - Install dependencies (streamdeck library, image libraries)
    - Implement device detection and initialization
    - Create button rendering (120×120 pixels)
    - Create touchscreen rendering (800×100 pixels)
    - Implement mode-specific section data
    - Add button, dial, and touch event handling
    - Integrate with main application
    - Set up Linux udev rules (if applicable)

16. **Stream Deck Tabs Architecture** (spec/19-streamdeck-tabs-architecture.md)
    - Design tabbed navigation system
    - Top row buttons select between 4 tabs
    - Tab 1: Light Control (existing functionality)
    - Tabs 2-4: Reserved for future features
    - Define tab-specific button and dial behavior

17. **Stream Deck Tabs Implementation** (spec/20-streamdeck-tabs-implementation-plan.md)
    - Add Tab type and state tracking
    - Update button handling for tab selection
    - Render tab buttons on top row
    - Make second row and touchscreen tab-aware
    - Add tab checks to dial handlers
    - Create placeholder content for future tabs

18. **Stream Deck Scenes Architecture** (spec/21-streamdeck-scenes-architecture.md)
    - Design Tab 2 for saving/recalling lighting scenes
    - Define database schema for scene storage
    - Plan 4 scene slots with save (dial click) and recall (button press)
    - Design touchscreen display showing scene status

19. **Stream Deck Scenes Implementation** (spec/22-streamdeck-scenes-implementation-plan.md)
    - Add scene tables to database schema
    - Implement scene storage operations (save/load/delete)
    - Add storage access to StreamDeckUI
    - Implement saveScene and recallScene methods
    - Update event handlers for Tab 2
    - Render scene buttons and touchscreen status

## Phase 9: Future Enhancements
20. **Additional UI Integration** (spec/08-future-ui-integration.md)
    - Add presets and scenes
    - WebSocket support for real-time updates
    - Animations and effects
    - Advanced Stream Deck features

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
- **Low Complexity:** Project setup, MQTT infrastructure, LED strip, video light, storage schema
- **Medium Complexity:** Main orchestration, testing, state storage operations, driver integration, TUI components (LED strip, video light), web backend (API, server)
- **High Complexity:** LED bar (due to complex message format), LED bar state storage, TUI LED bar component, web frontend (JavaScript, state management), Stream Deck+ interface (USB HID communication, image rendering, event handling)

## Dependencies
- Each driver depends on MQTT infrastructure
- Main orchestration depends on all drivers
- Testing can be done incrementally as each driver is completed
- State storage depends on drivers being implemented
- Driver state integration depends on storage layer
- TUI depends on drivers having getter methods for current state
- TUI depends on main orchestration and state storage being complete
- Web interface depends on drivers having getter methods for current state
- Web interface depends on main orchestration and state storage being complete
- Stream Deck interface depends on drivers having getter methods for current state
- Stream Deck interface depends on main orchestration and state storage being complete
- Web, TUI, and Stream Deck can be developed independently (parallel development possible)
- Stream Deck requires additional libraries (github.com/muesli/streamdeck, golang.org/x/image/...)
- Stream Deck requires hardware for full testing (can develop with mocks initially)
