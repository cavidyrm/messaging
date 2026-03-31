// cmd/migrate/main.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"messaging/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	mainMigrationsPath := filepath.Join(workDir, "migrations/main")
	eventMigrationsPath := filepath.Join(workDir, "migrations/events")

	if err := runMigrations(cfg.Database.DSN(), mainMigrationsPath); err != nil {
		log.Fatalf("Main DB migration failed: %v", err)
	}
	log.Println("Main DB migrations completed successfully")

	if err := runMigrations(cfg.EventDB.DSN(), eventMigrationsPath); err != nil {
		log.Fatalf("Event DB migration failed: %v", err)
	}
	log.Println("Event DB migrations completed successfully")
}

func runMigrations(dsn, path string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", path),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
