package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"portfolio/config"
	"portfolio/logger"
	migration "portfolio/migrations"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func NewConnection(cfg *config.DatabaseConfig, logger *logger.Logger) (*sql.DB, error) {
	dbDir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	dsn := buildDSN(cfg)

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if err := applyPragmas(db, cfg.Pragmas); err != nil {
		closeDB(db, logger)
		return nil, fmt.Errorf("failed to apply pragmas: %w", err)
	}

	if err := db.Ping(); err != nil {
		closeDB(db, logger)
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := applyMigrations(db, logger); err != nil {
		closeDB(db, logger)
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}

func buildDSN(cfg *config.DatabaseConfig) string {
	dsn := cfg.Path

	var options []string

	for key, value := range cfg.Options {
		options = append(options, fmt.Sprintf("%s=%s", key, value))
	}

	if _, exists := cfg.Options["_timeout"]; !exists {
		options = append(options, "_timeout=20000")
	}

	if _, exists := cfg.Options["_txlock"]; !exists {
		options = append(options, "_txlock=immediate")
	}

	if len(options) > 0 {
		dsn += "?" + strings.Join(options, "&")
	}

	return dsn
}

func applyPragmas(db *sql.DB, pragmas map[string]string) error {
	for pragma, value := range pragmas {
		query := fmt.Sprintf("PRAGMA %s = %s", pragma, value)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to apply pragma %s: %w", pragma, err)
		}
	}

	return verifyPragmas(db, pragmas)
}

func verifyPragmas(db *sql.DB, pragmas map[string]string) error {
	for pragma, expected := range pragmas {
		var current string
		query := fmt.Sprintf("PRAGMA %s", pragma)

		if err := db.QueryRow(query).Scan(&current); err != nil {
			return fmt.Errorf("failed to verify pragma %s: %w", pragma, err)
		}

		switch pragma {
		case "foreign_keys":
			if (expected == "ON" && current != "1") || (expected == "OFF" && current != "0") {
				return fmt.Errorf("pragma %s mismatch: expected %s, got %s", pragma, expected, current)
			}
		case "synchronous":
			// NORMAL = 1, FULL = 2, OFF = 0
			if expected == "NORMAL" && current != "1" {
				return fmt.Errorf("pragma %s mismatch: expected %s, got %s", pragma, expected, current)
			}
		default:
		}
	}

	return nil
}

func DiagnoseConnection(db *sql.DB, logger *logger.Logger) error {
	pragmas := []string{
		"foreign_keys", "journal_mode", "synchronous", "busy_timeout",
		"temp_store", "mmap_size", "journal_size_limit", "cache_size",
	}

	logger.Println("=== SQLite Configuration ===")
	for _, pragma := range pragmas {
		var value string
		query := fmt.Sprintf("PRAGMA %s", pragma)

		if err := db.QueryRow(query).Scan(&value); err != nil {
			logger.Printf("❌ %s: ERROR - %v\n", pragma, err)
		} else {
			logger.Printf("✅ %s: %s\n", pragma, value)
		}
	}
	logger.Println("=============================")

	return nil
}

func applyMigrations(db *sql.DB, logger *logger.Logger) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		schema_migration_id INTEGER PRIMARY KEY AUTOINCREMENT,
		schema_migration_filename TEXT NOT NULL UNIQUE,
		schema_migration_applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT schema_migration_filename FROM schema_migrations")
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Printf("Error closing rows: %v", err)
		}
	}()

	applied := make(map[string]bool)
	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return err
		}
		applied[filename] = true
	}

	migrationFiles, err := migration.GetMigrationFiles()

	if err != nil {
		return err
	}

	for _, migration := range migrationFiles {
		if applied[migration.Name] {
			continue
		}

		fmt.Println(string(migration.Content))

		if _, err := db.Exec(string(migration.Content)); err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO schema_migrations (schema_migration_filename) VALUES (?)", migration.Name); err != nil {
			return err
		}
		logger.Printf("Applied Migration: %s\n", migration.Name)
	}
	return nil
}

func closeDB(db *sql.DB, logger *logger.Logger) {
	defer func() {
		if err := db.Close(); err != nil {
			logger.Printf("db close: %v", err)
		}
	}()
}
