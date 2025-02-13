package testutils

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	testDB     *sql.DB
	testDBOnce sync.Once
)

// InitTestDB creates a test database and initializes its schema
func InitTestDB() *sql.DB {
	testDBOnce.Do(func() {
		// Create a temporary database file
		tempDir := os.TempDir()
		dbPath := filepath.Join(tempDir, "test_lang_portal.db")

		// Remove existing test database if it exists
		os.Remove(dbPath)

		// Open the database
		var err error
		testDB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
		if err != nil {
			log.Fatalf("Failed to open test database: %v", err)
		}

		// Set connection parameters
		testDB.SetMaxOpenConns(1)
		testDB.SetMaxIdleConns(1)

		// Define the schema with more explicit column definitions
		schema := `
		-- Words Table
		CREATE TABLE IF NOT EXISTS words (
			id INTEGER PRIMARY KEY,
			hindi TEXT NOT NULL,
			scrambled TEXT NOT NULL,
			hinglish TEXT NOT NULL,
			english TEXT NOT NULL,
			difficulty TEXT CHECK(difficulty IN ('easy', 'medium', 'hard')) DEFAULT 'medium',
			created_at DATETIME DEFAULT (datetime('now', 'localtime'))
		);

		-- Groups Table
		CREATE TABLE IF NOT EXISTS groups (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT DEFAULT '',
			created_at DATETIME DEFAULT (datetime('now', 'localtime'))
		);

		-- Word Groups Table (Many-to-Many Relationship)
		CREATE TABLE IF NOT EXISTS word_groups (
			word_id INTEGER,
			group_id INTEGER,
			created_at DATETIME DEFAULT (datetime('now', 'localtime')),
			PRIMARY KEY (word_id, group_id),
			FOREIGN KEY (word_id) REFERENCES words(id) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
		);

		-- Study Activities Table
		CREATE TABLE IF NOT EXISTS study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			description TEXT DEFAULT '',
			image TEXT DEFAULT '',
			score INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT (datetime('now', 'localtime'))
		);

		-- Sessions Table
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			activity_id INTEGER NOT NULL,
			group_id INTEGER,
			start_time DATETIME DEFAULT (datetime('now', 'localtime')),
			end_time DATETIME,
			score INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT (datetime('now', 'localtime')),
			FOREIGN KEY (activity_id) REFERENCES study_activities(id) ON DELETE CASCADE,
			FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
		);

		-- Session Activities Table
		CREATE TABLE IF NOT EXISTS session_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER NOT NULL,
			activity_id INTEGER NOT NULL,
			answer TEXT DEFAULT '',
			result TEXT DEFAULT '',
			score INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT (datetime('now', 'localtime')),
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
			FOREIGN KEY (activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
		);

		-- Indexes for performance
		CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);
		CREATE INDEX IF NOT EXISTS idx_sessions_activity ON sessions(activity_id);
		CREATE INDEX IF NOT EXISTS idx_session_activities_session ON session_activities(session_id);
		CREATE INDEX IF NOT EXISTS idx_session_activities_activity ON session_activities(activity_id);
		`

		// Execute schema creation
		_, err = testDB.Exec(schema)
		if err != nil {
			log.Fatalf("Failed to create test database schema: %v", err)
		}
	})

	return testDB
}

// CloseTestDB closes the test database connection
func CloseTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// ResetTestDB resets the test database by dropping and recreating all tables
func ResetTestDB(db *sql.DB) error {
	tables := []string{
		"session_activities", 
		"sessions", 
		"study_activities", 
		"word_groups", 
		"groups", 
		"words",
	}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table))
		if err != nil {
			return fmt.Errorf("failed to reset table %s: %w", table, err)
		}
	}

	return nil
}
