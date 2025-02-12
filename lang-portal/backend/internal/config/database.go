package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbDriver = "sqlite3"
	dbPath   = "./lang_portal.db"
)

// InitDatabase initializes and returns a database connection
func InitDatabase() (*sql.DB, error) {
	// Ensure the directory exists
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Open database connection
	db, err := sql.Open(dbDriver, dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established")
	return db, nil
}
