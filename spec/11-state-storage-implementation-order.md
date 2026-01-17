# State Storage Implementation Order

## Overview
Recommended implementation sequence for adding SQLite-based state persistence to the office lights system.

## Phase 1: Storage Layer Foundation

### Step 1: Set Up Dependencies
**Estimated complexity:** Low

1. Add SQLite driver dependency
   ```bash
   go get github.com/mattn/go-sqlite3
   ```
   Or for pure Go (no CGo):
   ```bash
   go get modernc.org/sqlite
   ```

2. Run `go mod tidy`

3. Verify the build still works

### Step 2: Create Storage Package Structure
**Estimated complexity:** Low

1. Create `storage/` directory
2. Create `storage/schema.go` with SQL schema definitions
3. Create `storage/database.go` with basic structure
4. Create `storage/interface.go` with `StateStore` interface
5. Create `storage/mock.go` with mock implementation for testing

### Step 3: Implement Schema Management
**Estimated complexity:** Medium

**File:** `storage/schema.go`

Tasks:
1. Define SQL CREATE TABLE statements as constants
2. Implement `InitSchema()` to create tables
3. Implement `InitDefaultData()` to insert default rows
4. Add schema existence checks
5. Write unit tests for schema creation

**Verification:**
- Run tests
- Manually inspect created database file with `sqlite3` CLI
- Verify tables exist with correct structure

### Step 4: Implement Database Connection
**Estimated complexity:** Low

**File:** `storage/database.go`

Tasks:
1. Implement `NewDatabase(path string)` constructor
2. Implement `Close()` method
3. Set SQLite pragmas (foreign keys, etc.)
4. Handle file creation
5. Write unit tests

**Verification:**
- Test database creation
- Test connection to existing database
- Test error handling for invalid paths

## Phase 2: Storage Operations

### Step 5: Implement LED Strip Storage
**Estimated complexity:** Low

**File:** `storage/database.go`

Tasks:
1. Implement `SaveLEDStripState(id, r, g, b int)`
   - Use INSERT OR REPLACE
   - Prepared statement
   - Transaction if needed
2. Implement `LoadLEDStripState(id int)`
   - Query with SELECT
   - Handle missing row (return defaults)
3. Write unit tests
   - Test save and load round-trip
   - Test missing data handling
   - Test invalid IDs

**Verification:**
- All tests pass
- Manually verify database contents

### Step 6: Implement Video Light Storage
**Estimated complexity:** Low

**File:** `storage/database.go`

Tasks:
1. Implement `SaveVideoLightState(id int, on bool, brightness int)`
   - Convert bool to int
   - Use INSERT OR REPLACE
2. Implement `LoadVideoLightState(id int)`
   - Convert int to bool
   - Handle missing row
3. Write unit tests

**Verification:**
- All tests pass
- Test both video light IDs (0 and 1)

### Step 7: Implement LED Bar Storage
**Estimated complexity:** High

**File:** `storage/database.go`

Tasks:
1. Implement `SaveLEDBarChannels(ledbarID int, channels []int)`
   - Validate 77 channels
   - Use transaction for atomic update
   - UPSERT each channel
   - Prepared statement for efficiency
2. Implement `LoadLEDBarChannels(ledbarID int)`
   - Query all 77 channels
   - Order by `channel_num`
   - Fill missing channels with 0
   - Return exactly 77 values
3. Write comprehensive unit tests
   - Test all 77 channels
   - Test partial data (missing channels)
   - Test transaction rollback on error
   - Test ordering

**Verification:**
- All tests pass
- Manually verify database has 77 rows per LED bar
- Check performance (should be <50ms for save)

## Phase 3: Driver Integration

### Step 8: Create Storage Mock
**Estimated complexity:** Low

**File:** `storage/mock.go`

Tasks:
1. Implement mock StateStore interface
2. Add call tracking for all methods
3. Add getter methods to verify calls
4. Write basic tests for mock itself

**Verification:**
- Mock compiles and implements interface
- Can track calls correctly

### Step 9: Update LED Strip Driver
**Estimated complexity:** Medium

**File:** `drivers/ledstrip/ledstrip.go`

Tasks:
1. Add `store` and `id` fields to struct
2. Create `NewLEDStripWithState()` constructor
3. Keep existing `NewLEDStrip()` for compatibility
4. Modify `Publish()` to save state
5. Add logging for storage errors
6. Update tests to use storage mock
7. Add new tests for state persistence

**Verification:**
- All existing tests still pass
- New tests verify state is saved
- No breaking changes to existing API

### Step 10: Update Video Light Driver
**Estimated complexity:** Medium

**File:** `drivers/videolight/videolight.go`

Tasks:
1. Add `store` field to struct
2. Create `NewVideoLightWithState()` constructor
3. Keep existing `NewVideoLight()` for compatibility
4. Modify `Publish()` to save state
5. Handle ID mapping (driver ID → database ID)
6. Update tests
7. Add state persistence tests

**Verification:**
- All tests pass
- Verify ID mapping is correct (1→0, 2→1)

### Step 11: Update LED Bar Driver
**Estimated complexity:** High

**File:** `drivers/ledbar/ledbar.go`

Tasks:
1. Add `store` field to struct
2. Create `NewLEDBarWithState()` constructor
3. Implement `loadFromChannels()` helper
4. Implement `getChannels()` helper
5. Keep existing `NewLEDBar()` for compatibility
6. Modify `Publish()` to save state
7. Update tests
8. Add comprehensive state persistence tests

**Verification:**
- All tests pass
- Test channel array round-trip
- Verify 77 channels are saved/loaded correctly

## Phase 4: Main Application Integration

### Step 12: Update Main.go - Database Setup
**Estimated complexity:** Medium

**File:** `main.go`

Tasks:
1. Add database connection logic
2. Get database path from environment variable
3. Initialize schema on startup
4. Add error handling for database failures
5. Close database in defer/shutdown

**Verification:**
- Application starts successfully
- Database file is created
- Tables are created
- Application shuts down cleanly

### Step 13: Update Main.go - Load State
**Estimated complexity:** Medium

**File:** `main.go`

Tasks:
1. Load LED strip state from database
2. Load LED bar state from database
3. Load video light states (both) from database
4. Handle missing/corrupt data gracefully
5. Log loaded states
6. Use loaded states in driver constructors

**Verification:**
- States are loaded correctly
- Defaults are used if database is empty
- Application handles corrupt data

### Step 14: Update Main.go - Initial Publish
**Estimated complexity:** Low

**File:** `main.go`

Tasks:
1. After creating drivers, publish initial state to MQTT
2. This ensures physical lights match stored state
3. Add logging for initial publish
4. Handle publish errors

**Verification:**
- Lights are set to stored state on startup
- MQTT messages are sent
- Logs confirm state restoration

### Step 15: Update Main.go - Shutdown Cleanup
**Estimated complexity:** Low

**File:** `main.go`

Tasks:
1. Ensure final state is saved before shutdown
2. Close database connection in correct order
3. Update shutdown sequence documentation

**Verification:**
- State is saved on graceful shutdown
- Database is closed cleanly
- No file locks or corruption

## Phase 5: Testing and Documentation

### Step 16: Integration Testing
**Estimated complexity:** Medium

Tasks:
1. Create integration test for full lifecycle:
   - Start application
   - Change lights
   - Shut down
   - Restart application
   - Verify state restored
2. Test with missing database
3. Test with corrupt database
4. Test concurrent operations (if applicable)

**Verification:**
- State persists across restarts
- All edge cases handled

### Step 17: Update Documentation
**Estimated complexity:** Low

Tasks:
1. Update `CONFIG.md` with database path configuration
2. Add database backup/restore instructions
3. Update `IMPLEMENTATION.md` with state storage details
4. Add troubleshooting section for database issues
5. Document database schema in separate file

**Verification:**
- All documentation is accurate
- Examples work as written

### Step 18: Performance Testing
**Estimated complexity:** Low

Tasks:
1. Measure database write latency
2. Verify LED bar saves (77 values) complete in <50ms
3. Test with SSD vs HDD (if applicable)
4. Document performance characteristics

**Verification:**
- Performance is acceptable
- No noticeable lag when changing lights

## Implementation Tips

### Development Workflow
1. **Test-Driven Development:** Write tests before implementation
2. **Incremental Testing:** Test each component before moving to next
3. **Database Inspection:** Use `sqlite3` CLI to verify database contents
4. **Backup Testing:** Keep test database files for debugging

### Testing Each Phase
```bash
# After each step, run:
go test ./...
go build
./office_lights  # Run briefly to test
```

### Database Inspection
```bash
# View database contents:
sqlite3 lights.sqlite3

# Useful commands:
.tables                    # List tables
.schema tablename         # Show table structure
SELECT * FROM ledstrips;  # Query data
```

### Debugging Tips
1. Add verbose logging during development
2. Use temporary database files for testing
3. Verify schema with `PRAGMA table_info(tablename)`
4. Test with empty/missing database file
5. Test with database from older version (future)

## Rollback Plan

If issues arise:
1. **Storage layer issues:** Revert storage package, drivers still work without it
2. **Driver integration issues:** Drivers fall back to no-op when store is nil
3. **Database corruption:** Delete `lights.sqlite3`, application creates new one
4. **Performance issues:** Add database tuning, investigate slow queries

## Validation Checklist

Before considering implementation complete:

### Storage Layer
- [ ] Database file is created automatically
- [ ] All tables created with correct schema
- [ ] Foreign key constraints enforced
- [ ] Unique constraints enforced
- [ ] Default values inserted
- [ ] Save operations work for all light types
- [ ] Load operations work for all light types
- [ ] Missing data handled gracefully
- [ ] All storage tests pass (>90% coverage)

### Driver Integration
- [ ] All drivers support state storage
- [ ] State saved after every publish
- [ ] Storage errors don't break MQTT operations
- [ ] Backward compatibility maintained
- [ ] All driver tests pass
- [ ] State persistence tests pass

### Main Application
- [ ] Database initialized on startup
- [ ] State loaded correctly
- [ ] Initial state published to MQTT
- [ ] State saved during operation
- [ ] Database closed on shutdown
- [ ] Environment variables work
- [ ] Graceful error handling

### End-to-End
- [ ] State persists across restarts
- [ ] All lights restore to previous state
- [ ] Database survives application crashes
- [ ] Performance is acceptable
- [ ] Documentation is complete
- [ ] No regressions in existing functionality

## Success Criteria
- All tests pass with >90% coverage
- State persists across application restarts
- Performance impact is negligible (<10ms per operation)
- No breaking changes to existing code
- Documentation is comprehensive
- Database is human-readable (can inspect with sqlite3)

## Estimated Total Time
- **Fast track (experienced):** 4-6 hours
- **Normal pace:** 8-12 hours
- **Learning while implementing:** 16-24 hours

## Next Steps After Completion
See `spec/08-future-ui-integration.md` for UI implementation that will benefit from state storage.
