package testutils

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const schema = `
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hindi TEXT NOT NULL,
    scrambled TEXT,
    hinglish TEXT,
    english TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
`

// CreateTestDB creates a temporary SQLite database for testing
func CreateTestDB() (*sql.DB, func(), error) {
	// Create a temporary directory for the test database
	tmpDir, err := os.MkdirTemp("", "langportal-test-*")
	if err != nil {
		return nil, nil, err
	}

	dbPath := filepath.Join(tmpDir, "test.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		return nil, nil, err
	}

	// Return the database and a cleanup function
	cleanup := func() {
		db.Close()
		os.RemoveAll(tmpDir)
	}

	// Create tables
	_, err = db.Exec(schema)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return db, cleanup, nil
}
