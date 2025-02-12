-- Lang Portal Database Schema
-- SQLite3 compatible schema with full Unicode support for Hindi

-- Ensure UTF-8 and Unicode support
PRAGMA encoding = 'UTF-8';
PRAGMA foreign_keys = ON;
PRAGMA recursive_triggers = ON;

-- Words Table with explicit TEXT support for multilingual content
CREATE TABLE IF NOT EXISTS words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    hindi TEXT NOT NULL,
    scrambled TEXT NOT NULL,
    hinglish TEXT NOT NULL,
    english TEXT NOT NULL,
    difficulty TEXT CHECK(difficulty IN ('easy', 'medium', 'hard')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create full-text search index for multilingual search
CREATE VIRTUAL TABLE IF NOT EXISTS words_fts USING fts5(
    hindi,
    hinglish,
    english,
    tokenize = 'unicode61'
);

-- Groups Table with Unicode support
CREATE TABLE IF NOT EXISTS groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Word Groups Table (Many-to-Many relationship between Words and Groups)
CREATE TABLE IF NOT EXISTS word_groups (
    word_id INTEGER,
    group_id INTEGER,
    PRIMARY KEY (word_id, group_id),
    FOREIGN KEY(word_id) REFERENCES words(id),
    FOREIGN KEY(group_id) REFERENCES groups(id)
);

-- Study Activities Table
CREATE TABLE IF NOT EXISTS study_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Sessions Table
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    activity_id INTEGER NOT NULL,
    group_id INTEGER,
    start_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    end_time DATETIME,
    score INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES study_activities(id) ON DELETE CASCADE,
    FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
);

-- Session Activities Table
CREATE TABLE IF NOT EXISTS session_activities (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session_id INTEGER NOT NULL,
    activity_id INTEGER NOT NULL,
    question TEXT NOT NULL,
    answer TEXT,
    result TEXT,
    score INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (activity_id) REFERENCES study_activities(id) ON DELETE CASCADE
);

-- Triggers to keep FTS index in sync
CREATE TRIGGER IF NOT EXISTS words_after_insert AFTER INSERT ON words BEGIN
    INSERT INTO words_fts(hindi, hinglish, english) 
    VALUES (NEW.hindi, NEW.hinglish, NEW.english);
END;

CREATE TRIGGER IF NOT EXISTS words_after_update AFTER UPDATE ON words BEGIN
    UPDATE words_fts 
    SET hindi = NEW.hindi, 
        hinglish = NEW.hinglish, 
        english = NEW.english 
    WHERE rowid = OLD.rowid;
END;

CREATE TRIGGER IF NOT EXISTS words_after_delete AFTER DELETE ON words BEGIN
    DELETE FROM words_fts WHERE rowid = OLD.rowid;
END;

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_words_hindi ON words(hindi);
CREATE INDEX IF NOT EXISTS idx_groups_name ON groups(name);
CREATE INDEX IF NOT EXISTS idx_sessions_activity ON sessions(activity_id);
CREATE INDEX IF NOT EXISTS idx_session_activities_session ON session_activities(session_id);
CREATE INDEX IF NOT EXISTS idx_session_activities_activity ON session_activities(activity_id);