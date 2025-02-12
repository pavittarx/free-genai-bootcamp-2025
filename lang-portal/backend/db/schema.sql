-- Lang Portal Database Schema
-- SQLite3 compatible schema with full Unicode support

-- Ensure UTF-8 and Unicode support
PRAGMA encoding = 'UTF-8';
PRAGMA foreign_keys = ON;
PRAGMA recursive_triggers = ON;

-- Words Table
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY,
    hindi TEXT NOT NULL,
    scrambled TEXT NOT NULL,
    hinglish TEXT NOT NULL,
    english TEXT NOT NULL,
    difficulty TEXT CHECK(difficulty IN ('easy', 'medium', 'hard')),
    created_at DATETIME DEFAULT (datetime('now', 'localtime'))
);

-- Create indexes for performance and text search
CREATE INDEX IF NOT EXISTS idx_word_difficulty ON words(difficulty);
CREATE INDEX IF NOT EXISTS idx_words_hindi ON words(hindi);
CREATE INDEX IF NOT EXISTS idx_words_hinglish ON words(hinglish);
CREATE INDEX IF NOT EXISTS idx_words_english ON words(english);

-- Groups Table
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
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

-- Create index for performance
CREATE INDEX IF NOT EXISTS idx_word_groups ON word_groups(word_id, group_id);

-- Study Activities Table
CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    image TEXT,
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
    answer TEXT,
    result TEXT,
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
