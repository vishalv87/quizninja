package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"quizninja-api/utils"

	"github.com/sirupsen/logrus"
)

// RunMigrations executes all SQL migration files in the migrations directory
func RunMigrations(db *sql.DB) error {
	// Create migrations table if it doesn't exist
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Run each migration
	for _, file := range migrationFiles {
		if err := runMigration(db, file); err != nil {
			return fmt.Errorf("failed to run migration %s: %w", file, err)
		}
	}

	utils.WithFields(logrus.Fields{
		"count": len(migrationFiles),
	}).Info("All migrations completed successfully")

	// Generate updated schema.sql file
	utils.Info("Updating schema.sql file")
	if err := UpdateSchemaAfterMigration(db); err != nil {
		utils.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Failed to update schema.sql")
		// Don't fail the migration if schema generation fails
	} else {
		utils.Info("Schema.sql updated successfully")
	}

	return nil
}

func createMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) UNIQUE NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.Exec(query)
	return err
}

func getMigrationFiles() ([]string, error) {
	var files []string

	// Try multiple possible migration directory paths
	migrationPaths := []string{
		"database/migrations",    // from project root
		"../database/migrations", // from tests directory
	}

	var migrationDir string
	var err error

	// Find the first path that exists
	for _, path := range migrationPaths {
		if _, statErr := os.Stat(path); statErr == nil {
			migrationDir = path
			break
		}
	}

	if migrationDir == "" {
		return nil, fmt.Errorf("migration directory not found in any of the expected locations: %v", migrationPaths)
	}

	err = filepath.WalkDir(migrationDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(path, ".sql") {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort files to ensure they run in order
	sort.Strings(files)
	return files, nil
}

func runMigration(db *sql.DB, filename string) error {
	// Check if migration has already been applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE filename = $1", filepath.Base(filename)).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		utils.WithFields(logrus.Fields{
			"migration": filepath.Base(filename),
		}).Debug("Migration already applied, skipping")
		return nil
	}

	// Read migration file
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration
	_, err = tx.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration as applied
	_, err = tx.Exec("INSERT INTO migrations (filename) VALUES ($1)", filepath.Base(filename))
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	utils.WithFields(logrus.Fields{
		"migration": filepath.Base(filename),
	}).Info("Successfully applied migration")
	return nil
}
