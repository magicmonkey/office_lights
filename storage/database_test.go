package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if db == nil {
		t.Fatal("Expected non-nil database")
	}

	// Verify file was created
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestInitSchema(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	err = db.InitSchema()
	if err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Verify tables exist
	tables := []string{"ledbars", "ledbars_leds", "ledstrips", "videolights"}
	for _, table := range tables {
		var name string
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
		err := db.db.QueryRow(query, table).Scan(&name)
		if err != nil {
			t.Errorf("Table %s does not exist: %v", table, err)
		}
	}
}

func TestInitDefaultData(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	if err := db.InitDefaultData(); err != nil {
		t.Fatalf("InitDefaultData failed: %v", err)
	}

	// Verify LED bar exists
	var count int
	err = db.db.QueryRow("SELECT COUNT(*) FROM ledbars WHERE id = 0").Scan(&count)
	if err != nil || count != 1 {
		t.Error("LED bar default data not inserted")
	}

	// Verify LED strip exists
	err = db.db.QueryRow("SELECT COUNT(*) FROM ledstrips WHERE id = 0").Scan(&count)
	if err != nil || count != 1 {
		t.Error("LED strip default data not inserted")
	}

	// Verify video lights exist
	err = db.db.QueryRow("SELECT COUNT(*) FROM videolights").Scan(&count)
	if err != nil || count != 2 {
		t.Errorf("Expected 2 video lights, got %d", count)
	}
}

func TestLEDStripStatePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	db.InitDefaultData()

	// Save state
	err = db.SaveLEDStripState(0, 100, 150, 200)
	if err != nil {
		t.Fatalf("Failed to save LED strip state: %v", err)
	}

	// Load state
	r, g, b, err := db.LoadLEDStripState(0)
	if err != nil {
		t.Fatalf("Failed to load LED strip state: %v", err)
	}

	if r != 100 || g != 150 || b != 200 {
		t.Errorf("Expected RGB (100,150,200), got (%d,%d,%d)", r, g, b)
	}
}

func TestLEDStripStateUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	db.InitDefaultData()

	// Save initial state
	db.SaveLEDStripState(0, 100, 100, 100)

	// Update state
	err = db.SaveLEDStripState(0, 255, 0, 0)
	if err != nil {
		t.Fatalf("Failed to update LED strip state: %v", err)
	}

	// Load and verify updated state
	r, g, b, _ := db.LoadLEDStripState(0)
	if r != 255 || g != 0 || b != 0 {
		t.Errorf("State not updated correctly: got (%d,%d,%d)", r, g, b)
	}
}

func TestVideoLightStatePersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	db.InitDefaultData()

	// Save video light 0 state
	err = db.SaveVideoLightState(0, true, 75)
	if err != nil {
		t.Fatalf("Failed to save video light state: %v", err)
	}

	// Load state
	on, brightness, err := db.LoadVideoLightState(0)
	if err != nil {
		t.Fatalf("Failed to load video light state: %v", err)
	}

	if !on {
		t.Error("Expected light to be on")
	}
	if brightness != 75 {
		t.Errorf("Expected brightness 75, got %d", brightness)
	}
}

func TestVideoLightBooleanConversion(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	db.InitDefaultData()

	// Test true
	db.SaveVideoLightState(0, true, 50)
	on, _, _ := db.LoadVideoLightState(0)
	if !on {
		t.Error("Expected true, got false")
	}

	// Test false
	db.SaveVideoLightState(0, false, 50)
	on, _, _ = db.LoadVideoLightState(0)
	if on {
		t.Error("Expected false, got true")
	}
}

func TestLEDBarChannelsPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	db.InitDefaultData()

	// Create test channels (77 values)
	channels := make([]int, 77)
	for i := 0; i < 77; i++ {
		channels[i] = i + 1 // Values 1-77
	}

	// Save channels
	err = db.SaveLEDBarChannels(0, channels)
	if err != nil {
		t.Fatalf("Failed to save LED bar channels: %v", err)
	}

	// Load channels
	loadedChannels, err := db.LoadLEDBarChannels(0)
	if err != nil {
		t.Fatalf("Failed to load LED bar channels: %v", err)
	}

	// Verify
	if len(loadedChannels) != 77 {
		t.Fatalf("Expected 77 channels, got %d", len(loadedChannels))
	}

	for i := 0; i < 77; i++ {
		if loadedChannels[i] != i+1 {
			t.Errorf("Channel %d: expected %d, got %d", i, i+1, loadedChannels[i])
		}
	}
}

func TestLEDBarChannelsUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	db.InitDefaultData()

	// Save initial channels
	channels1 := make([]int, 77)
	for i := range channels1 {
		channels1[i] = 100
	}
	db.SaveLEDBarChannels(0, channels1)

	// Update channels
	channels2 := make([]int, 77)
	for i := range channels2 {
		channels2[i] = 200
	}
	db.SaveLEDBarChannels(0, channels2)

	// Load and verify
	loaded, _ := db.LoadLEDBarChannels(0)
	for i, val := range loaded {
		if val != 200 {
			t.Errorf("Channel %d not updated: expected 200, got %d", i, val)
		}
	}
}

func TestLEDBarChannelsInvalidLength(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	db.InitDefaultData()

	// Try to save wrong number of channels
	wrongChannels := make([]int, 50)
	err = db.SaveLEDBarChannels(0, wrongChannels)
	if err == nil {
		t.Error("Expected error for wrong channel count")
	}
}

func TestLoadNonExistentData(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	db.InitSchema()
	// Don't initialize default data

	// Load LED strip (should return defaults)
	r, g, b, err := db.LoadLEDStripState(0)
	if err != nil {
		t.Errorf("Loading non-existent LED strip should not error: %v", err)
	}
	if r != 0 || g != 0 || b != 0 {
		t.Errorf("Expected default (0,0,0), got (%d,%d,%d)", r, g, b)
	}

	// Load video light (should return defaults)
	on, brightness, err := db.LoadVideoLightState(0)
	if err != nil {
		t.Errorf("Loading non-existent video light should not error: %v", err)
	}
	if on || brightness != 0 {
		t.Errorf("Expected default (false, 0), got (%v, %d)", on, brightness)
	}

	// Load LED bar (should return 77 zeros)
	channels, err := db.LoadLEDBarChannels(0)
	if err != nil {
		t.Errorf("Loading non-existent LED bar should not error: %v", err)
	}
	if len(channels) != 77 {
		t.Fatalf("Expected 77 channels, got %d", len(channels))
	}
	for i, val := range channels {
		if val != 0 {
			t.Errorf("Channel %d: expected 0, got %d", i, val)
		}
	}
}

func TestDatabaseClose(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}

	// Try to close again (should not panic)
	err = db.Close()
	if err != nil {
		t.Errorf("Second close failed: %v", err)
	}
}
