package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func InitDB() {
	var err error
	// Register sqlite driver is done by the import
	db, err = sql.Open("sqlite", "notes.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Create table
	schema := `
	CREATE TABLE IF NOT EXISTS notes (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		tags TEXT DEFAULT '[]',
		updated_at TEXT NOT NULL,
		created_at TEXT NOT NULL
	);
	`
	if _, err = db.Exec(schema); err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Optimization for concurrency
	db.SetMaxOpenConns(1) // SQLite single writer, though WAL mode allows more.
	// But modernc/sqlite + pure go, let's play safe with 1 to avoid locking issues in this simple app unless we enable WAL.

	// Enable WAL mode for better concurrency if we wanted, but sticking to simple default for "small assessment" is fine.
	// Actually, setting busy timeout is good practice.
	_, _ = db.Exec("PRAGMA busy_timeout = 5000;")
}
