#!/bin/bash

# Set error handling
set -e

# Database path
DB_PATH="../lang_portal.db"

# Migration script path
MIGRATION_SCRIPT="../db/migrations/001_initial_schema.sql"

# Check if database exists
if [ -f "$DB_PATH" ]; then
    echo "Database already exists. Skipping initialization."
    exit 0
fi

# Create database and apply migrations
sqlite3 "$DB_PATH" < "$MIGRATION_SCRIPT"

echo "Database initialized successfully."
