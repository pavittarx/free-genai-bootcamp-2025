# Database Generation and Seeding

## Overview
This package provides functionality to generate and seed the SQLite database for the Language Portal project.

## Features
- Create a new SQLite database
- Apply database schema
- Seed data from CSV files
- Support for multiple tables (words, groups, word_groups)

## Prerequisites
- SQLite3 installed
- CSV seed files in `../seeds/` directory
- Executable permissions on `init.sh`

## Usage

### Generate Database
```bash
# From the backend directory
./db/init.sh
```

### Features
- Removes existing database
- Creates new SQLite database
- Seeds data from CSV files
- Verifies data after seeding

## Data Sources
- `../seeds/groups.csv`: Group definitions
- `../seeds/words.csv`: Word entries
- `../seeds/word_groups.csv`: Word-to-Group mappings

## Seed Data Files
- `groups.csv`: Group definitions
- `words.csv`: Word entries
- `word_groups.csv`: Word-to-Group mappings

## Notes
- Existing database will be overwritten
- Requires CSV files to be present in the seeds directory
- Supports easy, medium, and hard difficulty levels

## Troubleshooting
- Ensure SQLite3 is installed
- Check CSV files for correct format
- Verify file permissions