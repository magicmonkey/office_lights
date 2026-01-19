package storage

const (
	// SQL schema for the lights database
	schemaLEDBars = `
CREATE TABLE IF NOT EXISTS ledbars (
    id INTEGER PRIMARY KEY
);`

	schemaLEDBarsLEDs = `
CREATE TABLE IF NOT EXISTS ledbars_leds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ledbar_id INTEGER NOT NULL,
    channel_num INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK(value >= 0 AND value <= 255),
    FOREIGN KEY (ledbar_id) REFERENCES ledbars(id) ON DELETE CASCADE,
    UNIQUE(ledbar_id, channel_num)
);`

	schemaLEDStrips = `
CREATE TABLE IF NOT EXISTS ledstrips (
    id INTEGER PRIMARY KEY,
    red INTEGER NOT NULL DEFAULT 0 CHECK(red >= 0 AND red <= 255),
    green INTEGER NOT NULL DEFAULT 0 CHECK(green >= 0 AND green <= 255),
    blue INTEGER NOT NULL DEFAULT 0 CHECK(blue >= 0 AND blue <= 255)
);`

	schemaVideoLights = `
CREATE TABLE IF NOT EXISTS videolights (
    id INTEGER PRIMARY KEY,
    "on" INTEGER NOT NULL DEFAULT 0 CHECK("on" IN (0, 1)),
    brightness INTEGER NOT NULL DEFAULT 0 CHECK(brightness >= 0 AND brightness <= 100)
);`

	schemaIndex = `
CREATE INDEX IF NOT EXISTS idx_ledbars_leds_lookup
ON ledbars_leds(ledbar_id, channel_num);`

	// Scene tables for saving/recalling light presets
	schemaScenes = `
CREATE TABLE IF NOT EXISTS scenes (
    id INTEGER PRIMARY KEY
);`

	schemaScenesLEDBarsLEDs = `
CREATE TABLE IF NOT EXISTS scenes_ledbars_leds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    ledbar_id INTEGER NOT NULL,
    channel_num INTEGER NOT NULL,
    value INTEGER NOT NULL CHECK(value >= 0 AND value <= 255),
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE
);`

	schemaScenesLEDStrips = `
CREATE TABLE IF NOT EXISTS scenes_ledstrips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    red INTEGER NOT NULL CHECK(red >= 0 AND red <= 255),
    green INTEGER NOT NULL CHECK(green >= 0 AND green <= 255),
    blue INTEGER NOT NULL CHECK(blue >= 0 AND blue <= 255),
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE
);`

	schemaScenesVideoLights = `
CREATE TABLE IF NOT EXISTS scenes_videolights (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scene_id INTEGER NOT NULL,
    videolight_id INTEGER NOT NULL,
    on_state INTEGER NOT NULL CHECK(on_state IN (0, 1)),
    brightness INTEGER NOT NULL CHECK(brightness >= 0 AND brightness <= 100),
    FOREIGN KEY (scene_id) REFERENCES scenes(id) ON DELETE CASCADE
);`

	schemaScenesIndex = `
CREATE INDEX IF NOT EXISTS idx_scenes_ledbars_leds_lookup
ON scenes_ledbars_leds(scene_id, ledbar_id, channel_num);`

	// Default data initialization
	initLEDBars = `INSERT OR IGNORE INTO ledbars (id) VALUES (0);`

	initLEDStrips = `INSERT OR IGNORE INTO ledstrips (id, red, green, blue) VALUES (0, 0, 0, 0);`

	initVideoLights = `
INSERT OR IGNORE INTO videolights (id, "on", brightness) VALUES (0, 0, 0);
INSERT OR IGNORE INTO videolights (id, "on", brightness) VALUES (1, 0, 0);`

	// Initialize 4 empty scene slots (IDs 0-3)
	initScenes = `
INSERT OR IGNORE INTO scenes (id) VALUES (0);
INSERT OR IGNORE INTO scenes (id) VALUES (1);
INSERT OR IGNORE INTO scenes (id) VALUES (2);
INSERT OR IGNORE INTO scenes (id) VALUES (3);`
)

// allSchemas returns all CREATE TABLE statements in order
func allSchemas() []string {
	return []string{
		schemaLEDBars,
		schemaLEDBarsLEDs,
		schemaLEDStrips,
		schemaVideoLights,
		schemaIndex,
		schemaScenes,
		schemaScenesLEDBarsLEDs,
		schemaScenesLEDStrips,
		schemaScenesVideoLights,
		schemaScenesIndex,
	}
}

// allInitData returns all default data INSERT statements in order
func allInitData() []string {
	return []string{
		initLEDBars,
		initLEDStrips,
		initVideoLights,
		initScenes,
	}
}
