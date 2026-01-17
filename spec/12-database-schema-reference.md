# Database Schema Reference

## Overview
This document provides a complete reference for the SQLite database schema used for state persistence.

## Database File
- **Filename:** `lights.sqlite3`
- **Location:** Current working directory
- **Format:** SQLite3
- **Expected Size:** <100KB

## Schema Definitions

### `ledbars` Table

Stores LED bar instances.

```sql
CREATE TABLE IF NOT EXISTS ledbars (
    id INTEGER PRIMARY KEY
);
```

**Columns:**
- `id` (INTEGER, PRIMARY KEY): LED bar identifier
  - Always 0 for the single LED bar in the system

**Initial Data:**
```sql
INSERT OR IGNORE INTO ledbars (id) VALUES (0);
```

**Example Data:**
```
id
--
0
```

---

### `ledbars_leds` Table

Stores individual LED channel values for LED bars.

```sql
CREATE TABLE IF NOT EXISTS ledbars_leds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ledbar_id INTEGER NOT NULL,
    channel_num INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK(value >= 0 AND value <= 255),
    FOREIGN KEY (ledbar_id) REFERENCES ledbars(id) ON DELETE CASCADE,
    UNIQUE(ledbar_id, channel_num)
);
```

**Columns:**
- `id` (INTEGER, PRIMARY KEY, AUTOINCREMENT): Unique row identifier
- `ledbar_id` (INTEGER, NOT NULL, FK): References `ledbars.id`
- `channel_num` (INTEGER, NOT NULL): Channel number (0-76)
- `value` (INTEGER, NOT NULL): LED value (0-255)

**Constraints:**
- `CHECK(value >= 0 AND value <= 255)`: Ensures valid LED values
- `UNIQUE(ledbar_id, channel_num)`: One value per channel per bar
- `FOREIGN KEY`: Cascade delete when LED bar is deleted

**Channel Mapping:**
LED bar has 77 channels (0-76) representing:
- Channels 0-23: First 6 RGBW LEDs (R,G,B,W × 6)
- Channels 24-36: First 13 white LEDs
- Channels 37-39: 3 ignored values (always 0)
- Channels 40-63: Second 6 RGBW LEDs (R,G,B,W × 6)
- Channels 64-76: Second 13 white LEDs

**Example Data:**
```
id  | ledbar_id | channel_num | value
----|-----------|-------------|------
1   | 0         | 0           | 255   (First RGBW LED, Red)
2   | 0         | 1           | 0     (First RGBW LED, Green)
3   | 0         | 2           | 0     (First RGBW LED, Blue)
4   | 0         | 3           | 100   (First RGBW LED, White)
... | ...       | ...         | ...
77  | 0         | 76          | 50    (Last white LED)
```

---

### `ledstrips` Table

Stores LED strip RGB state.

```sql
CREATE TABLE IF NOT EXISTS ledstrips (
    id INTEGER PRIMARY KEY,
    red INTEGER NOT NULL DEFAULT 0 CHECK(red >= 0 AND red <= 255),
    green INTEGER NOT NULL DEFAULT 0 CHECK(green >= 0 AND green <= 255),
    blue INTEGER NOT NULL DEFAULT 0 CHECK(blue >= 0 AND blue <= 255)
);
```

**Columns:**
- `id` (INTEGER, PRIMARY KEY): LED strip identifier
  - Always 0 for the single LED strip in the system
- `red` (INTEGER, NOT NULL, DEFAULT 0): Red value (0-255)
- `green` (INTEGER, NOT NULL, DEFAULT 0): Green value (0-255)
- `blue` (INTEGER, NOT NULL, DEFAULT 0): Blue value (0-255)

**Constraints:**
- `CHECK(red >= 0 AND red <= 255)`: Valid red value
- `CHECK(green >= 0 AND green <= 255)`: Valid green value
- `CHECK(blue >= 0 AND blue <= 255)`: Valid blue value

**Initial Data:**
```sql
INSERT OR IGNORE INTO ledstrips (id, red, green, blue) VALUES (0, 0, 0, 0);
```

**Example Data:**
```
id | red | green | blue
---|-----|-------|-----
0  | 255 | 200   | 150   (Warm white)
```

---

### `videolights` Table

Stores video light states.

```sql
CREATE TABLE IF NOT EXISTS videolights (
    id INTEGER PRIMARY KEY,
    on INTEGER NOT NULL DEFAULT 0 CHECK(on IN (0, 1)),
    brightness INTEGER NOT NULL DEFAULT 0 CHECK(brightness >= 0 AND brightness <= 100)
);
```

**Columns:**
- `id` (INTEGER, PRIMARY KEY): Video light identifier
  - 0 for video light 1
  - 1 for video light 2
- `on` (INTEGER, NOT NULL, DEFAULT 0): On/off state (0=off, 1=on)
- `brightness` (INTEGER, NOT NULL, DEFAULT 0): Brightness (0-100)

**Constraints:**
- `CHECK(on IN (0, 1))`: Boolean stored as 0 or 1
- `CHECK(brightness >= 0 AND brightness <= 100)`: Valid brightness

**Initial Data:**
```sql
INSERT OR IGNORE INTO videolights (id, on, brightness) VALUES (0, 0, 0);
INSERT OR IGNORE INTO videolights (id, on, brightness) VALUES (1, 0, 0);
```

**Example Data:**
```
id | on | brightness
---|----|-----------
0  | 1  | 75        (Video light 1: on at 75%)
1  | 0  | 50        (Video light 2: off, last brightness 50%)
```

---

## ID Mappings

### LED Bar
- **Database ID:** 0
- **Driver Instance:** Single instance, ID 0
- **MQTT Topic:** `kevinoffice/ledbar/0`

### LED Strip
- **Database ID:** 0
- **Driver Instance:** Single instance, ID 0
- **MQTT Topic:** `kevinoffice/ledstrip/sequence`

### Video Lights
- **Database ID:** 0 → **Driver ID:** 1 → **MQTT Topic:** `kevinoffice/videolight/1/command/light:0`
- **Database ID:** 1 → **Driver ID:** 2 → **MQTT Topic:** `kevinoffice/videolight/2/command/light:0`

**Important:** Video lights use 0-based IDs in database but 1-based IDs in drivers. The mapping is:
- Database ID = Driver ID - 1

---

## Indexes

Consider adding indexes for better query performance:

```sql
-- Index for LED bar channel lookups
CREATE INDEX IF NOT EXISTS idx_ledbars_leds_lookup
ON ledbars_leds(ledbar_id, channel_num);

-- Not strictly necessary given small dataset, but good practice
```

---

## Pragmas

Recommended SQLite pragmas for the application:

```sql
-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

-- Use Write-Ahead Logging for better concurrency
PRAGMA journal_mode = WAL;

-- Synchronous mode (NORMAL is good balance of safety and speed)
PRAGMA synchronous = NORMAL;
```

---

## Database Queries

### Common Queries

#### Load LED Strip State
```sql
SELECT red, green, blue
FROM ledstrips
WHERE id = 0;
```

#### Save LED Strip State
```sql
INSERT OR REPLACE INTO ledstrips (id, red, green, blue)
VALUES (0, ?, ?, ?);
```

#### Load LED Bar Channels
```sql
SELECT channel_num, value
FROM ledbars_leds
WHERE ledbar_id = 0
ORDER BY channel_num;
```

#### Save Single LED Bar Channel
```sql
INSERT OR REPLACE INTO ledbars_leds (ledbar_id, channel_num, value)
VALUES (0, ?, ?);
```

#### Save All LED Bar Channels (Transaction)
```sql
BEGIN TRANSACTION;

DELETE FROM ledbars_leds WHERE ledbar_id = 0;

INSERT INTO ledbars_leds (ledbar_id, channel_num, value) VALUES
    (0, 0, ?), (0, 1, ?), (0, 2, ?), ..., (0, 76, ?);

COMMIT;
```

Or using UPSERT in a loop:
```sql
BEGIN TRANSACTION;

-- For each channel 0-76:
INSERT OR REPLACE INTO ledbars_leds (ledbar_id, channel_num, value)
VALUES (0, ?, ?);

COMMIT;
```

#### Load Video Light State
```sql
SELECT on, brightness
FROM videolights
WHERE id = ?;
```

#### Save Video Light State
```sql
INSERT OR REPLACE INTO videolights (id, on, brightness)
VALUES (?, ?, ?);
```

---

## Database Initialization

Complete initialization script:

```sql
-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Create tables
CREATE TABLE IF NOT EXISTS ledbars (
    id INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS ledbars_leds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ledbar_id INTEGER NOT NULL,
    channel_num INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK(value >= 0 AND value <= 255),
    FOREIGN KEY (ledbar_id) REFERENCES ledbars(id) ON DELETE CASCADE,
    UNIQUE(ledbar_id, channel_num)
);

CREATE TABLE IF NOT EXISTS ledstrips (
    id INTEGER PRIMARY KEY,
    red INTEGER NOT NULL DEFAULT 0 CHECK(red >= 0 AND red <= 255),
    green INTEGER NOT NULL DEFAULT 0 CHECK(green >= 0 AND green <= 255),
    blue INTEGER NOT NULL DEFAULT 0 CHECK(blue >= 0 AND blue <= 255)
);

CREATE TABLE IF NOT EXISTS videolights (
    id INTEGER PRIMARY KEY,
    on INTEGER NOT NULL DEFAULT 0 CHECK(on IN (0, 1)),
    brightness INTEGER NOT NULL DEFAULT 0 CHECK(brightness >= 0 AND brightness <= 100)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_ledbars_leds_lookup
ON ledbars_leds(ledbar_id, channel_num);

-- Insert default data
INSERT OR IGNORE INTO ledbars (id) VALUES (0);
INSERT OR IGNORE INTO ledstrips (id, red, green, blue) VALUES (0, 0, 0, 0);
INSERT OR IGNORE INTO videolights (id, on, brightness) VALUES (0, 0, 0);
INSERT OR IGNORE INTO videolights (id, on, brightness) VALUES (1, 0, 0);
```

---

## Database Inspection

### Using SQLite CLI

```bash
# Open database
sqlite3 lights.sqlite3

# List tables
.tables

# Show schema
.schema

# Show specific table schema
.schema ledstrips

# Query data
SELECT * FROM ledstrips;
SELECT * FROM videolights;
SELECT * FROM ledbars_leds ORDER BY channel_num;

# Show table info
PRAGMA table_info(ledbars_leds);

# Check foreign keys
PRAGMA foreign_key_list(ledbars_leds);

# Compact database
VACUUM;

# Exit
.quit
```

### Useful Queries for Debugging

```sql
-- Count LED bar channels (should be 77)
SELECT COUNT(*) FROM ledbars_leds WHERE ledbar_id = 0;

-- Find missing channels
WITH RECURSIVE
  cnt(x) AS (
    SELECT 0
    UNION ALL
    SELECT x+1 FROM cnt
    LIMIT 77
  )
SELECT x FROM cnt
WHERE x NOT IN (SELECT channel_num FROM ledbars_leds WHERE ledbar_id = 0);

-- Show all non-zero channels
SELECT channel_num, value
FROM ledbars_leds
WHERE ledbar_id = 0 AND value > 0
ORDER BY channel_num;

-- Check for invalid values
SELECT * FROM ledbars_leds WHERE value < 0 OR value > 255;
SELECT * FROM ledstrips WHERE red < 0 OR red > 255 OR green < 0 OR green > 255 OR blue < 0 OR blue > 255;
SELECT * FROM videolights WHERE on NOT IN (0, 1) OR brightness < 0 OR brightness > 100;
```

---

## Backup and Recovery

### Backup
```bash
# Simple file copy (when application is stopped)
cp lights.sqlite3 lights.sqlite3.backup

# Backup while application is running (uses SQLite backup API)
sqlite3 lights.sqlite3 ".backup lights.sqlite3.backup"

# Export to SQL
sqlite3 lights.sqlite3 .dump > lights.sql
```

### Restore
```bash
# From backup file
cp lights.sqlite3.backup lights.sqlite3

# From SQL dump
sqlite3 lights.sqlite3.new < lights.sql
mv lights.sqlite3.new lights.sqlite3
```

### Reset Database
```bash
# Delete database file - application will recreate
rm lights.sqlite3
```

---

## Troubleshooting

### Database Locked Error
- Another process has the database open
- Check for zombie processes
- Wait and retry (implement retry logic)
- Ensure WAL mode is enabled

### Corrupted Database
```bash
# Check integrity
sqlite3 lights.sqlite3 "PRAGMA integrity_check;"

# If corrupted, try to recover
sqlite3 lights.sqlite3 ".recover" > recovered.sql
sqlite3 lights.sqlite3.new < recovered.sql

# Or start fresh (lose data)
rm lights.sqlite3
```

### Missing Data
- Check if schema was initialized
- Verify default data was inserted
- Look for application errors in logs

### Performance Issues
- Enable WAL mode
- Add indexes
- Use transactions for bulk operations
- Check disk I/O

---

## Future Enhancements

### Schema Versioning
Consider adding a version table for schema migrations:

```sql
CREATE TABLE schema_version (
    version INTEGER PRIMARY KEY,
    applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO schema_version (version) VALUES (1);
```

### Audit Logging
Track when states change:

```sql
CREATE TABLE state_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    light_type TEXT NOT NULL,
    light_id INTEGER NOT NULL,
    changed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    old_value TEXT,
    new_value TEXT
);
```

### Scenes/Presets
Store named scenes:

```sql
CREATE TABLE scenes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE scene_states (
    scene_id INTEGER NOT NULL,
    light_type TEXT NOT NULL,
    light_id INTEGER NOT NULL,
    state_data TEXT NOT NULL,
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE
);
```

---

## Testing

### Test Data Setup

```sql
-- Setup test data for LED strip (warm white)
INSERT OR REPLACE INTO ledstrips (id, red, green, blue)
VALUES (0, 255, 200, 150);

-- Setup test data for video lights
INSERT OR REPLACE INTO videolights (id, on, brightness)
VALUES (0, 1, 75);
INSERT OR REPLACE INTO videolights (id, on, brightness)
VALUES (1, 1, 50);

-- Setup test data for LED bar (first LED red)
INSERT OR REPLACE INTO ledbars_leds (ledbar_id, channel_num, value)
VALUES (0, 0, 255), (0, 1, 0), (0, 2, 0), (0, 3, 0);
```

### Clean Test Database

```bash
# Create fresh test database
rm -f test.sqlite3
sqlite3 test.sqlite3 < schema.sql
```
