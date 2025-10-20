package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// New creates a new database connection
func New(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(1) // SQLite works best with a single connection
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Run migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	migrationFile := "migrations/00001_initial_schema.sql"

	// Read migration file
	content, err := os.ReadFile(migrationFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Check if tables already exist
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='companies'").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check for existing tables: %w", err)
	}

	// If tables already exist, skip migration
	if count > 0 {
		return nil
	}

	// Extract SQL from goose annotations
	lines := string(content)

	// Find the "Up" section markers
	upStart := "-- +goose Up\n-- +goose StatementBegin\n"
	upEnd := "\n-- +goose StatementEnd"

	startIdx := len(upStart)
	endIdx := len(lines)

	// Find the end marker position
	if idx := findString(lines, upEnd); idx != -1 {
		endIdx = idx
	}

	// Extract SQL content
	sqlContent := lines[startIdx:endIdx]

	// Execute migration
	if _, err := db.Exec(sqlContent); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	return nil
}

func findString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
