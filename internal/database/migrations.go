package database

import (
	"context"
	"fmt"
)

func (db *Database) RunMigrations() error {
	ctx := context.Background()

	// Start a transaction
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	// Defer rollback in case of panic or early return
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Create migrations table if it doesn't exist
	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INT PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// List of migrations (add new steps here)
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS devices (
			id SERIAL PRIMARY KEY,
			fuseId BIGINT UNIQUE NOT NULL,
			name VARCHAR(255),
			description TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS sensor_data (
			id SERIAL PRIMARY KEY,
			device_id INT REFERENCES devices(id) ON DELETE CASCADE,
			topic_id INT NOT NULL,
			payload_version INT NOT NULL,
			payload VARCHAR(128) NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS location VARCHAR(255);`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS device_type VARCHAR(100);`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS wifi_strength INT;`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS battery_percent INT;`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS type VARCHAR(50);`,
		`ALTER TABLE devices ADD COLUMN IF NOT EXISTS last_seen TIMESTAMPTZ DEFAULT NOW();`,
	}

	// Apply migrations sequentially
	for i, query := range migrations {
		version := i + 1
		var exists bool
		err = tx.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)", version,
		).Scan(&exists)
		if err != nil {
			return fmt.Errorf("failed to check migration %d: %w", version, err)
		}

		if !exists {
			fmt.Printf("Applying migration %d...\n", version)
			_, err = tx.Exec(ctx, query)
			if err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", version, err)
			}

			_, err = tx.Exec(ctx,
				"INSERT INTO schema_migrations (version) VALUES ($1)", version,
			)
			if err != nil {
				return fmt.Errorf("failed to record migration %d: %w", version, err)
			}
		}
	}

	// Commit transaction if all migrations succeed
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}

	fmt.Println("All migrations applied successfully.")
	return nil
}
