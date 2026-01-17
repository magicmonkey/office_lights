# State Storage with SQLite

## Overview
Implement persistent state storage using SQLite3 to preserve light states across application restarts.

## Database Design

### Database File
- **Location:** `./lights.sqlite3` (current directory)
- **Format:** SQLite3 database

### Table Structures

#### `ledbars` Table
```sql
CREATE TABLE ledbars (
    id INTEGER PRIMARY KEY
);
```
- Stores LED bar instances
- Single row with `id = 0` (hard-coded)

#### `ledbars_leds` Table
```sql
CREATE TABLE ledbars_leds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ledbar_id INTEGER NOT NULL,
    channel_num INTEGER NOT NULL,
    value INTEGER NOT NULL,
    FOREIGN KEY (ledbar_id) REFERENCES ledbars(id),
    UNIQUE(ledbar_id, channel_num)
);
```
- Stores individual LED channel values for LED bars
- `ledbar_id`: Reference to `ledbars.id` (always 0)
- `channel_num`: Channel number (0-76, representing the 77 values in the CSV)
- `value`: LED brightness/color value (0-255)

#### `ledstrips` Table
```sql
CREATE TABLE ledstrips (
    id INTEGER PRIMARY KEY,
    red INTEGER NOT NULL DEFAULT 0,
    green INTEGER NOT NULL DEFAULT 0,
    blue INTEGER NOT NULL DEFAULT 0
);
```
- Stores LED strip RGB state
- Single row with `id = 0` (hard-coded)

#### `videolights` Table
```sql
CREATE TABLE videolights (
    id INTEGER PRIMARY KEY,
    on INTEGER NOT NULL DEFAULT 0,
    brightness INTEGER NOT NULL DEFAULT 0
);
```
- Stores video light states
- Two rows: `id = 0` and `id = 1` (hard-coded for the two video lights)
- `on`: Boolean stored as INTEGER (0 = false, 1 = true)
- `brightness`: 0-100

## Implementation Tasks

### 1. Create Storage Package
- Create `storage/` directory
- Create `storage/database.go` for database operations
- Create `storage/schema.go` for schema definitions

### 2. Database Connection Management
- `NewDatabase(path string) (*Database, error)`: Open/create database
- `Close() error`: Close database connection
- Implement connection pooling if needed
- Handle SQLite-specific pragmas (e.g., `PRAGMA foreign_keys = ON`)

### 3. Schema Initialization
- `InitSchema() error`: Create tables if they don't exist
- `InitDefaultData() error`: Insert default rows for:
  - LED bar (id=0)
  - LED strip (id=0)
  - Video lights (id=0, id=1)
- Check if schema already exists before creating

### 4. LED Strip Storage Operations
- `SaveLEDStripState(id int, r, g, b int) error`
  - Update `ledstrips` table
  - Use UPSERT (INSERT OR REPLACE) for simplicity
- `LoadLEDStripState(id int) (r, g, b int, error)`
  - Query `ledstrips` table
  - Return RGB values

### 5. LED Bar Storage Operations
- `SaveLEDBarChannels(ledbarID int, channels []int) error`
  - Update all 77 channels in `ledbars_leds` table
  - Use transaction for atomic update
  - UPSERT each channel value
- `LoadLEDBarChannels(ledbarID int) ([]int, error)`
  - Query all channels for the LED bar, ordered by `channel_num`
  - Return array of 77 values
  - If fewer than 77 rows exist, fill missing with 0

### 6. Video Light Storage Operations
- `SaveVideoLightState(id int, on bool, brightness int) error`
  - Update `videolights` table
  - Convert bool to int (0/1)
  - Use UPSERT
- `LoadVideoLightState(id int) (on bool, brightness int, error)`
  - Query `videolights` table
  - Convert int to bool for `on` field

### 7. Transaction Support
- Implement transaction wrapper for bulk operations
- Ensure LED bar channel updates are atomic

### 8. Error Handling
- Handle database locked errors (common with SQLite)
- Handle missing database file (create new)
- Handle corrupted database (log error, start fresh?)
- Validate data ranges when loading from database

## Integration with Drivers

### Modify LED Strip Driver
- Add `SaveState()` method that calls storage layer
- Call `SaveState()` after every `Publish()`
- Accept initial state in constructor from storage

### Modify LED Bar Driver
- Add `SaveState()` method that calls storage layer
- Call `SaveState()` after every `Publish()`
- Accept initial state in constructor from storage

### Modify Video Light Driver
- Add `SaveState()` method that calls storage layer
- Call `SaveState()` after every `Publish()`
- Accept initial state in constructor from storage

## Main.go Integration

### Startup Sequence
1. Open database connection
2. Initialize schema and default data
3. Load state for all lights:
   - LED strip (id=0)
   - LED bar (id=0)
   - Video lights (id=0, id=1)
4. Create driver instances with loaded state
5. Publish initial state to MQTT (so physical lights match stored state)

### Shutdown Sequence
1. Turn off all lights (existing behavior)
2. Save final state to database
3. Close database connection
4. Disconnect MQTT

## Dependencies

### SQLite Driver
- Use `github.com/mattn/go-sqlite3` (CGo-based, most popular)
- Alternative: `modernc.org/sqlite` (pure Go, no CGo)
- Add to `go.mod` via `go get`

### Database Access Pattern
- Use `database/sql` standard library
- Prepared statements for repeated queries
- Transactions for bulk operations

## Testing Strategy

### Unit Tests
- Create `storage/database_test.go`
- Use temporary database files for testing
- Test schema creation
- Test each CRUD operation
- Test transaction rollback
- Test handling of missing/corrupted data

### Integration Tests
- Test full startup sequence (load → create drivers → publish)
- Test state persistence (save → close → reopen → load → verify)
- Test concurrent access (if applicable)

### Test Database Cleanup
- Use `t.TempDir()` for test database files
- Ensure each test gets a fresh database

## Performance Considerations

### Write Performance
- LED bar has 77 values: Use single transaction for all updates
- Consider batching writes if performance is an issue
- SQLite write latency is typically <10ms for small operations

### Read Performance
- State is loaded only once at startup
- No performance concerns for read operations

### Database Size
- Expected size: <100KB (very small dataset)
- No archival/cleanup needed

## Configuration

### Optional Environment Variables
- `DB_PATH`: Override default database path
  - Default: `./lights.sqlite3`
- `DB_RESET`: Reset database on startup (for testing)
  - Set to `1` to recreate database

## Migration Path

### For Existing Deployments
Since this is a new feature:
1. First run will create new database
2. Default values will be all lights off
3. First MQTT publish will save actual state

### Schema Versioning (Future)
- Consider adding `schema_version` table for future migrations
- Not required for initial implementation

## Example Usage

```go
// In main.go
db, err := storage.NewDatabase("lights.sqlite3")
if err != nil {
    log.Fatal(err)
}
defer db.Close()

if err := db.InitSchema(); err != nil {
    log.Fatal(err)
}

// Load LED strip state
r, g, b, err := db.LoadLEDStripState(0)
if err != nil {
    log.Printf("Error loading LED strip state, using defaults: %v", err)
    r, g, b = 0, 0, 0
}

// Create driver with loaded state
strip := ledstrip.NewLEDStripWithState(mqttClient, topic, r, g, b)

// Driver saves state after every change
strip.SetColor(255, 100, 50) // Automatically saves to DB
```

## Success Criteria
- Database file is created on first run
- All tables are created with correct schema
- State is loaded on startup and applied to drivers
- State is saved after every MQTT publish
- Application survives restart with state preserved
- Tests achieve >90% coverage for storage package
- No data loss under normal operation

## Security Considerations
- Database file has no sensitive data
- File permissions: default user permissions acceptable
- No SQL injection risk (using prepared statements)
- No authentication/encryption needed

## Backup and Recovery
- Database file can be backed up by copying `lights.sqlite3`
- To reset state: delete database file and restart application
- Consider documenting backup/restore procedure in CONFIG.md
