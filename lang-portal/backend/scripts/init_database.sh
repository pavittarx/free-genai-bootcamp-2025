#!/bin/bash

# Set the path to the project directory
PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

# Database file path
DB_FILE="${PROJECT_DIR}/lang-portal.db"

# Current timestamp
TIMESTAMP="$(date '+%Y-%m-%d %H:%M:%S')"

# Remove existing database if it exists
rm -f "${DB_FILE}"

# Initialize SQLite database
sqlite3 "${DB_FILE}" << EOF
-- Create schema
.read ${PROJECT_DIR}/db/schema.sql

-- Import groups with timestamp
.mode csv
.import --skip 1 ${PROJECT_DIR}/db/seeds/groups.csv groups

-- Update groups created_at
UPDATE groups SET created_at = '${TIMESTAMP}' WHERE created_at IS NULL;

-- Import words with timestamp
.import --skip 1 ${PROJECT_DIR}/db/seeds/words.csv words

-- Update words created_at
UPDATE words SET created_at = '${TIMESTAMP}' WHERE created_at IS NULL;

-- Import word groups with timestamp
.import --skip 1 ${PROJECT_DIR}/db/seeds/word_groups.csv word_groups

-- Update word_groups created_at
UPDATE word_groups SET created_at = '${TIMESTAMP}' WHERE created_at IS NULL;

-- Populate full-text search index
INSERT INTO words_fts (hindi, hinglish, english)
SELECT hindi, hinglish, english FROM words;

-- Verify imports
SELECT 'Groups count: ' || COUNT(*) FROM groups;
SELECT 'Words count: ' || COUNT(*) FROM words;
SELECT 'Word Groups count: ' || COUNT(*) FROM word_groups;
SELECT 'First group created_at: ' || created_at FROM groups LIMIT 1;
SELECT 'First word created_at: ' || created_at FROM words LIMIT 1;
SELECT 'First word_group created_at: ' || created_at FROM word_groups LIMIT 1;
EOF

# Check the exit status
if [ $? -eq 0 ]; then
    echo "Database initialized successfully!"
else
    echo "Error initializing database."
    exit 1
fi
