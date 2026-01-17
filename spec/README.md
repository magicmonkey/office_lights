# Office Lights Specifications

This directory contains detailed specification documents for implementing the office lights control system.

## Quick Start

**Already implemented (Steps 1-10):**
- ✅ Project setup and MQTT infrastructure
- ✅ All three light drivers (LED strip, LED bar, video lights)
- ✅ Main orchestration and testing
- ✅ State storage with SQLite
- ✅ Driver state integration

**Next to implement (Steps 11-12):**
- 📋 Text user interface (TUI)

**Future:**
- 🔮 Additional UI integration (Stream Deck, web interface)

## Specification Files

### Core Implementation (Completed ✅)

#### 00. [Implementation Order](00-implementation-order.md)
Master guide for implementation sequence and dependencies.

#### 01. [Project Setup](01-project-setup.md)
Go module initialization, project structure, and dependencies.
- **Status:** ✅ Complete
- **Complexity:** Low

#### 02. [MQTT Infrastructure](02-mqtt-infrastructure.md)
MQTT client wrapper, connection management, and topic definitions.
- **Status:** ✅ Complete
- **Complexity:** Low

#### 03. [LED Strip Driver](03-ledstrip-driver.md)
RGB LED strip control with JSON message formatting.
- **Status:** ✅ Complete
- **Complexity:** Low
- **Coverage:** 94.7%

#### 04. [LED Bar Driver](04-ledbar-driver.md)
RGBW LED bar with 77-channel CSV message format.
- **Status:** ✅ Complete
- **Complexity:** High
- **Coverage:** 93.2%

#### 05. [Video Light Driver](05-videolight-driver.md)
Video light on/off and brightness control.
- **Status:** ✅ Complete
- **Complexity:** Low
- **Coverage:** 96.4%

#### 06. [Main Orchestration](06-main-orchestration.md)
Application startup, driver management, and shutdown.
- **Status:** ✅ Complete
- **Complexity:** Medium

#### 07. [Testing Strategy](07-testing-strategy.md)
Unit tests, integration tests, and mocking.
- **Status:** ✅ Complete
- **Complexity:** Medium

---

### State Persistence (Completed ✅)

#### 09. [State Storage](09-state-storage.md)
SQLite database for persistent light state.
- **Status:** ✅ Complete
- **Complexity:** Medium
- **Coverage:** 58.9%

**What it covers:**
- Database schema design
- SQLite integration
- Load/save operations for all light types
- Error handling and recovery

#### 10. [Driver State Integration](10-driver-state-integration.md)
Integrate state storage into existing drivers.
- **Status:** ✅ Complete
- **Complexity:** Medium

**What it covers:**
- Adding StateStore interface to drivers
- Automatic state persistence on publish
- Loading initial state on startup
- Backward compatibility

#### 11. [State Storage Implementation Order](11-state-storage-implementation-order.md)
Step-by-step guide for implementing state persistence.
- **Type:** Implementation guide
- **Status:** ✅ Followed and completed

**What it covers:**
- 18 detailed implementation steps
- Testing strategy for each phase
- Rollback plans
- Validation checklists

#### 12. [Database Schema Reference](12-database-schema-reference.md)
Complete database schema documentation.
- **Type:** Reference document
- **Status:** ✅ Complete

**What it covers:**
- Table definitions
- SQL queries
- Database inspection
- Backup/recovery procedures

---

### Text User Interface (To Do 📋)

#### 13. [TUI Architecture](13-tui-architecture.md)
Terminal-based user interface design.
- **Status:** 📋 Not started
- **Complexity:** Medium

**What it covers:**
- Bubbletea framework selection
- Component architecture (4 light sections)
- Keyboard navigation design
- State synchronization with drivers
- Layout and styling approach

#### 14. [TUI Implementation Plan](14-tui-implementation-plan.md)
Step-by-step guide for implementing the TUI.
- **Status:** 📋 Not started
- **Complexity:** Medium
- **Estimated time:** 8-12 hours

**What it covers:**
- 8 implementation phases
- Component models for each light type
- Update and view functions
- Key binding implementation
- Integration with main application
- Testing strategy

---

### Future Enhancements (🔮)

#### 08. [UI Integration](08-future-ui-integration.md)
Additional user interface options beyond TUI.
- **Status:** 🔮 Future work
- **Complexity:** High

**What it covers:**
- Stream Deck integration
- Web interface
- API layer design
- Scene presets

---

## Implementation Progress

### Phase 1: Foundation ✅
- [x] Project setup
- [x] MQTT infrastructure

### Phase 2: Drivers ✅
- [x] LED strip driver
- [x] Video light driver
- [x] LED bar driver

### Phase 3: Integration ✅
- [x] Main orchestration
- [x] Testing

### Phase 4: State Persistence ✅
- [x] Storage infrastructure (spec 09)
- [x] Driver integration (spec 10)

### Phase 5: Text User Interface 📋
- [ ] TUI architecture (spec 13)
- [ ] TUI implementation (spec 14)

### Phase 6: Future 🔮
- [ ] Additional UI integration (spec 08)

---

## Getting Started with TUI

If you're ready to implement the text user interface, follow this sequence:

1. **Read the architecture:** Start with [13-tui-architecture.md](13-tui-architecture.md)
2. **Follow the implementation plan:** Use [14-tui-implementation-plan.md](14-tui-implementation-plan.md)
3. **Test incrementally:** Build and test each component as you go

### Quick Implementation Checklist

- [ ] Install Bubbletea framework: `go get github.com/charmbracelet/bubbletea@latest`
- [ ] Install Bubbles components: `go get github.com/charmbracelet/bubbles@latest`
- [ ] Install Lipgloss styling: `go get github.com/charmbracelet/lipgloss@latest`
- [ ] Create `tui/` package structure
- [ ] Add getter methods to drivers (R(), G(), B(), IsOn(), Brightness())
- [ ] Implement LED strip component
- [ ] Implement video light components
- [ ] Implement LED bar component
- [ ] Create root model and update function
- [ ] Implement view rendering with layout
- [ ] Add keyboard navigation
- [ ] Integrate with main.go (TUI mode)
- [ ] Test all controls and navigation
- [ ] Verify MQTT publishing and database saves

---

## Project Structure

```
office_lights/
├── main.go                  # Application entry point
├── mqtt/                    # MQTT client wrapper
│   ├── client.go
│   ├── topics.go
│   └── mock.go
├── drivers/                 # Light drivers
│   ├── ledstrip/
│   ├── ledbar/
│   └── videolight/
├── storage/                 # State persistence
│   ├── database.go
│   ├── schema.go
│   ├── interface.go
│   └── mock.go
├── tui/                     # Text user interface (to be created)
│   ├── tui.go
│   ├── model.go
│   ├── update.go
│   ├── view.go
│   ├── keys.go
│   ├── ledstrip.go
│   ├── ledbar.go
│   ├── videolight.go
│   ├── styles.go
│   └── messages.go
├── spec/                    # This directory
│   └── *.md
└── lights.sqlite3          # State database (created at runtime)
```

---

## Documentation

### User Documentation
- [README.md](../README.md) - Project overview
- [CONFIG.md](../CONFIG.md) - Configuration guide
- [USAGE.md](../USAGE.md) - Programming examples
- [IMPLEMENTATION.md](../IMPLEMENTATION.md) - Implementation summary

### Development Documentation
- This directory (`spec/`) - Detailed specifications
- Test files (`*_test.go`) - Test examples and patterns
- Code comments - Inline documentation

---

## Key Concepts

### Light Types
1. **LED Strip:** RGB color control (1 instance)
2. **LED Bar:** RGBW + white LEDs with 77 channels (1 instance)
3. **Video Lights:** On/off and brightness (2 instances)

### MQTT Topics
- `kevinoffice/ledstrip/sequence` - LED strip
- `kevinoffice/ledbar/0` - LED bar
- `kevinoffice/videolight/1/command/light:0` - Video light 1
- `kevinoffice/videolight/2/command/light:0` - Video light 2

### Message Formats
- **LED Strip:** JSON `{"sequence":"fill", "data":{"r":100,"g":150,"b":200}}`
- **LED Bar:** CSV with 77 comma-separated values
- **Video Lights:** Text `set,<on>,<brightness>`

### State Storage
- **Database:** SQLite3 (`lights.sqlite3`)
- **Tables:** `ledstrips`, `ledbars`, `ledbars_leds`, `videolights`
- **Behavior:** Load on startup, save on every change

---

## Testing

### Current Test Coverage
- LED Strip: 94.7%
- LED Bar: 93.2%
- Video Lights: 96.4%

### Test Commands
```bash
# Run all tests
go test ./...

# Run with coverage
go test -v -cover ./drivers/...

# Run specific package
go test ./drivers/ledstrip/
```

---

## Questions?

For questions about:
- **Implementation details:** Check the relevant spec file
- **Database schema:** See [12-database-schema-reference.md](12-database-schema-reference.md)
- **Step-by-step guide:** See [11-state-storage-implementation-order.md](11-state-storage-implementation-order.md)
- **Testing approaches:** See [07-testing-strategy.md](07-testing-strategy.md)

---

## Contributing

When adding new features:
1. Create a spec file in this directory
2. Update this README with the new spec
3. Update [00-implementation-order.md](00-implementation-order.md) if it affects implementation order
4. Follow the existing spec format for consistency

---

## Version History

- **v1.0** - Initial specifications (specs 01-07)
- **v1.1** - Added state storage specifications (specs 09-12)
- **v1.2** - Added TUI specifications (specs 13-14)
- **v2.0** - Will include additional UI integration (spec 08)
