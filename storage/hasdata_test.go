package storage

import (
	"path/filepath"
	"testing"
)

func TestHasData(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Check data before initialization - should be false
	hasData, err := db.HasData()
	if err != nil {
		t.Fatalf("HasData failed: %v", err)
	}
	if hasData {
		t.Error("Expected HasData to return false for empty database")
	}

	// Initialize default data
	if err := db.InitDefaultData(); err != nil {
		t.Fatalf("InitDefaultData failed: %v", err)
	}

	// Check data after initialization - should be true
	hasData, err = db.HasData()
	if err != nil {
		t.Fatalf("HasData failed: %v", err)
	}
	if !hasData {
		t.Error("Expected HasData to return true after initialization")
	}
}

func TestHasDataWithPartialData(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Should be false initially
	hasData, _ := db.HasData()
	if hasData {
		t.Error("Expected no data initially")
	}

	// Add only LED strip data
	if err := db.SaveLEDStripState(0, 100, 150, 200); err != nil {
		t.Fatalf("Failed to save LED strip state: %v", err)
	}

	// Should now have data
	hasData, err = db.HasData()
	if err != nil {
		t.Fatalf("HasData failed: %v", err)
	}
	if !hasData {
		t.Error("Expected HasData to return true after saving LED strip")
	}
}

func TestHasDataWithVideoLightsOnly(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Add only video light data
	if err := db.SaveVideoLightState(0, true, 75); err != nil {
		t.Fatalf("Failed to save video light state: %v", err)
	}

	// Should have data
	hasData, err := db.HasData()
	if err != nil {
		t.Fatalf("HasData failed: %v", err)
	}
	if !hasData {
		t.Error("Expected HasData to return true after saving video light")
	}
}

func TestHasDataWithLEDBarOnly(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	db, err := NewDatabase(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	if err := db.InitSchema(); err != nil {
		t.Fatalf("InitSchema failed: %v", err)
	}

	// Add LED bar instance
	_, err = db.db.Exec("INSERT INTO ledbars (id) VALUES (0)")
	if err != nil {
		t.Fatalf("Failed to insert LED bar: %v", err)
	}

	// Should have data
	hasData, err := db.HasData()
	if err != nil {
		t.Fatalf("HasData failed: %v", err)
	}
	if !hasData {
		t.Error("Expected HasData to return true after adding LED bar")
	}
}
