package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // SQLite driver
)

// Database provides persistent storage for light states using SQLite
type Database struct {
	db   *sql.DB
	path string
}

// NewDatabase creates a new database connection
func NewDatabase(path string) (*Database, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Set pragmas for better performance and reliability
	pragmas := []string{
		"PRAGMA foreign_keys = ON",
		"PRAGMA journal_mode = WAL",
		"PRAGMA synchronous = NORMAL",
	}

	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to set pragma: %w", err)
		}
	}

	return &Database{
		db:   db,
		path: path,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// InitSchema creates all tables and indexes if they don't exist
func (d *Database) InitSchema() error {
	log.Println("Storage: Initializing database schema...")

	// Execute all schema statements
	for _, schema := range allSchemas() {
		if _, err := d.db.Exec(schema); err != nil {
			return fmt.Errorf("failed to create schema: %w", err)
		}
	}

	log.Println("Storage: Schema initialized successfully")
	return nil
}

// InitDefaultData inserts default data if it doesn't exist
func (d *Database) InitDefaultData() error {
	log.Println("Storage: Initializing default data...")

	// Execute all init data statements
	for _, stmt := range allInitData() {
		if _, err := d.db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to init data: %w", err)
		}
	}

	log.Println("Storage: Default data initialized successfully")
	return nil
}

// SaveLEDStripState saves the RGB state for an LED strip
func (d *Database) SaveLEDStripState(id int, r, g, b int) error {
	query := `INSERT OR REPLACE INTO ledstrips (id, red, green, blue) VALUES (?, ?, ?, ?)`

	_, err := d.db.Exec(query, id, r, g, b)
	if err != nil {
		return fmt.Errorf("failed to save LED strip state: %w", err)
	}

	return nil
}

// LoadLEDStripState loads the RGB state for an LED strip
func (d *Database) LoadLEDStripState(id int) (r, g, b int, err error) {
	query := `SELECT red, green, blue FROM ledstrips WHERE id = ?`

	err = d.db.QueryRow(query, id).Scan(&r, &g, &b)
	if err == sql.ErrNoRows {
		// No data found, return defaults
		return 0, 0, 0, nil
	}
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to load LED strip state: %w", err)
	}

	return r, g, b, nil
}

// SaveVideoLightState saves the on/off and brightness state for a video light
func (d *Database) SaveVideoLightState(id int, on bool, brightness int) error {
	onInt := 0
	if on {
		onInt = 1
	}

	query := `INSERT OR REPLACE INTO videolights (id, "on", brightness) VALUES (?, ?, ?)`

	_, err := d.db.Exec(query, id, onInt, brightness)
	if err != nil {
		return fmt.Errorf("failed to save video light state: %w", err)
	}

	return nil
}

// LoadVideoLightState loads the on/off and brightness state for a video light
func (d *Database) LoadVideoLightState(id int) (on bool, brightness int, err error) {
	query := `SELECT "on", brightness FROM videolights WHERE id = ?`

	var onInt int
	err = d.db.QueryRow(query, id).Scan(&onInt, &brightness)
	if err == sql.ErrNoRows {
		// No data found, return defaults
		return false, 0, nil
	}
	if err != nil {
		return false, 0, fmt.Errorf("failed to load video light state: %w", err)
	}

	on = onInt == 1
	return on, brightness, nil
}

// SaveLEDBarChannels saves all 77 channel values for an LED bar
func (d *Database) SaveLEDBarChannels(ledbarID int, channels []int) error {
	if len(channels) != 77 {
		return fmt.Errorf("expected 77 channels, got %d", len(channels))
	}

	// Use a transaction for atomic update
	tx, err := d.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare statement for efficiency
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO ledbars_leds (ledbar_id, channel_num, value) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Insert all channels
	for i, value := range channels {
		if _, err := stmt.Exec(ledbarID, i, value); err != nil {
			return fmt.Errorf("failed to save channel %d: %w", i, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// LoadLEDBarChannels loads all 77 channel values for an LED bar
func (d *Database) LoadLEDBarChannels(ledbarID int) ([]int, error) {
	query := `SELECT channel_num, value FROM ledbars_leds WHERE ledbar_id = ? ORDER BY channel_num`

	rows, err := d.db.Query(query, ledbarID)
	if err != nil {
		return nil, fmt.Errorf("failed to load LED bar channels: %w", err)
	}
	defer rows.Close()

	// Create result array with all 77 channels initialized to 0
	channels := make([]int, 77)

	// Fill in values from database
	for rows.Next() {
		var channelNum, value int
		if err := rows.Scan(&channelNum, &value); err != nil {
			return nil, fmt.Errorf("failed to scan channel: %w", err)
		}

		if channelNum >= 0 && channelNum < 77 {
			channels[channelNum] = value
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating channels: %w", err)
	}

	return channels, nil
}
