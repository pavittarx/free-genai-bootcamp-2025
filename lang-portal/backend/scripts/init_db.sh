#!/bin/bash

# Database Initialization Script for Language Portal

# Directories
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
DB_PATH="$PROJECT_ROOT/lang_portal.db"
MIGRATIONS_PATH="$PROJECT_ROOT/db/migrations"
SEEDS_PATH="$PROJECT_ROOT/db/seeds"

# Colors for logging
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Remove existing database
if [ -f "$DB_PATH" ]; then
    rm "$DB_PATH"
    echo -e "${RED}Removed existing database${NC}"
fi

# Create new SQLite database
sqlite3 "$DB_PATH" << EOF
-- Apply migrations
.read "$MIGRATIONS_PATH/001_initial.sql"

-- Import groups
.mode csv
.import "$SEEDS_PATH/groups.csv" groups

-- Import words
.import "$SEEDS_PATH/words.csv" words

-- Import word_groups
.import "$SEEDS_PATH/word_groups.csv" word_groups

-- Verify imports
SELECT COUNT(*) as group_count FROM groups;
SELECT COUNT(*) as word_count FROM words;
SELECT COUNT(*) as word_group_count FROM word_groups;
EOF

# Check the exit status of the SQLite command
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Database initialized successfully at $DB_PATH${NC}"
    echo -e "${GREEN}Total groups: $(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM groups;")${NC}"
    echo -e "${GREEN}Total words: $(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM words;")${NC}"
else
    echo -e "${RED}Database initialization failed${NC}"
    exit 1
fi
