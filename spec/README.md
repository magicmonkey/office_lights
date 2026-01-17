# Office Lights Specifications

This directory contains detailed specification documents for implementing the office lights control system.

## Quick Start

**Already implemented (Steps 1-7):**
- ✅ Project setup and MQTT infrastructure
- ✅ All three light drivers (LED strip, LED bar, video lights)
- ✅ Main orchestration and testing

**Next to implement (Steps 8-9):**
- 📋 State storage with SQLite
- 📋 Driver state integration

**Future:**
- 🔮 UI integration

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

### State Persistence (To Do 📋)

#### 09. [State Storage](09-state-storage.md)
SQLite database for persistent light state.
- **Status:** 📋 Not started
- **Complexity:** Medium
- **Estimated time:** 4-8 hours

**What it covers:**
- Database schema design
- SQLite integration
- Load/save operations for all light types
- Error handling and recovery

#### 10. [Driver State Integration](10-driver-state-integration.md)
Integrate state storage into existing drivers.
- **Status:** 📋 Not started
- **Complexity:** Medium
- **Estimated time:** 3-6 hours

**What it covers:**
- Adding StateStore interface to drivers
- Automatic state persistence on publish
- Loading initial state on startup
- Backward compatibility

#### 11. [State Storage Implementation Order](11-state-storage-implementation-order.md)
Step-by-step guide for implementing state persistence.
- **Type:** Implementation guide
- **Complexity:** Detailed walkthrough

**What it covers:**
- 18 detailed implementation steps
- Testing strategy for each phase
- Rollback plans
- Validation checklists

#### 12. [Database Schema Reference](12-database-schema-reference.md)
Complete database schema documentation.
- **Type:** Reference document
- **Purpose:** Schema details, queries, troubleshooting

**What it covers:**
- Table definitions
- SQL queries
- Database inspection
- Backup/recovery procedures

---

### Future Enhancements (🔮)

#### 08. [UI Integration](08-future-ui-integration.md)
User interface planning and Stream Deck integration.
- **Status:** 🔮 Future work
- **Complexity:** High

**What it covers:**
- UI framework options
- API layer design
- Control layouts
- Scene presets
- Stream Deck integration

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

### Phase 4: State Persistence 📋
- [ ] Storage infrastructure (spec 09)
- [ ] Driver integration (spec 10)

### Phase 5: Future 🔮
- [ ] UI integration (spec 08)

---

## Getting Started with State Storage

If you're ready to implement state persistence, follow this sequence:

1. **Read the overview:** Start with [09-state-storage.md](09-state-storage.md)
2. **Understand driver changes:** Review [10-driver-state-integration.md](10-driver-state-integration.md)
3. **Follow the steps:** Use [11-state-storage-implementation-order.md](11-state-storage-implementation-order.md)
4. **Reference the schema:** Keep [12-database-schema-reference.md](12-database-schema-reference.md) handy

### Quick Implementation Checklist

- [ ] Install SQLite driver: `go get github.com/mattn/go-sqlite3`
- [ ] Create `storage/` package
- [ ] Implement database schema and initialization
- [ ] Implement save/load operations
- [ ] Create storage mock for testing
- [ ] Update LED strip driver with state storage
- [ ] Update video light driver with state storage
- [ ] Update LED bar driver with state storage
- [ ] Update main.go to load/save state
- [ ] Test state persistence across restarts
- [ ] Update documentation

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
├── storage/                 # State persistence (to be created)
│   ├── database.go
│   ├── schema.go
│   ├── interface.go
│   └── mock.go
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
- **v2.0** - Will include UI integration (spec 08)
