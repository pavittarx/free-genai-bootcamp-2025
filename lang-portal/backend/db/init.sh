#!/bin/bash

# Database Generation and Seeding Script

# Fail on any error
set -e

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Paths
PROJECT_ROOT="$(dirname "$0")/.."
DB_FILE="${PROJECT_ROOT}/lang_portal.db"
SEEDS_DIR="${PROJECT_ROOT}/seeds"
SCHEMA_FILE="${PROJECT_ROOT}/db/schema.sql"

# Function to log messages
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check prerequisites
check_prerequisites() {
    # Check for sqlite3
    if ! command -v sqlite3 &> /dev/null; then
        log_error "sqlite3 is not installed. Please install it first."
        exit 1
    fi

    # Check seed files
    for file in groups words word_groups; do
        if [ ! -f "${SEEDS_DIR}/${file}.csv" ]; then
            log_error "Missing seed file: ${SEEDS_DIR}/${file}.csv"
            exit 1
        fi
    done

    # Check schema file
    if [ ! -f "${SCHEMA_FILE}" ]; then
        log_error "Missing schema file: ${SCHEMA_FILE}"
        exit 1
    fi
}

# Remove existing database
clean_database() {
    if [ -f "${DB_FILE}" ]; then
        rm "${DB_FILE}"
        log_info "Removed existing database"
    fi
}

# Create database and apply schema
create_database() {
    log_info "Creating new database"
    sqlite3 "${DB_FILE}" < "${SCHEMA_FILE}"
}

# Seed groups table
seed_groups() {
    log_info "Seeding groups table"
    sqlite3 "${DB_FILE}" <<EOF
.mode csv
.import ${SEEDS_DIR}/groups.csv groups
EOF
}

# Seed words table
seed_words() {
    log_info "Seeding words table"
    sqlite3 "${DB_FILE}" <<EOF
.mode csv
.import ${SEEDS_DIR}/words.csv words
EOF
}

# Seed word_groups table
seed_word_groups() {
    log_info "Seeding word_groups table"
    sqlite3 "${DB_FILE}" <<EOF
.mode csv
.import ${SEEDS_DIR}/word_groups.csv word_groups
EOF
}

# Verify seeded data
verify_data() {
    log_info "Verifying seeded data"
    
    groups_count=$(sqlite3 "${DB_FILE}" "SELECT COUNT(*) FROM groups;")
    words_count=$(sqlite3 "${DB_FILE}" "SELECT COUNT(*) FROM words;")
    word_groups_count=$(sqlite3 "${DB_FILE}" "SELECT COUNT(*) FROM word_groups;")

    log_info "Groups: ${groups_count}"
    log_info "Words: ${words_count}"
    log_info "Word-Group Associations: ${word_groups_count}"
}

# Main script execution
main() {
    log_info "Starting database generation process"

    # Run checks
    check_prerequisites

    # Clean and recreate database
    clean_database
    create_database

    # Seed data
    seed_groups
    seed_words
    seed_word_groups

    # Verify data
    verify_data

    log_info "Database generation completed successfully!"
}

# Run the main function
main

exit 0
