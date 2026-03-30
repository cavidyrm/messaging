package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"messaging/config"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	if err := runMigrations(cfg.Database, "migrations"); err != nil {
		log.Fatal("Failed to run migrations on main DB:", err)
	}

	if err := runMigrations(cfg.EventDB, "migrations"); err != nil {
		log.Fatal("Failed to run migrations on event DB:", err)
	}

	log.Println("Migrations completed successfully")
}

func runMigrations(dbCfg config.DatabaseConfig, migrationsDir string) error {
	db, err := sql.Open("postgres", dbCfg.DSN())
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		log.Printf("Running migration: %s\n", file)
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute %s: %w", file, err)
		}
	}

	return nil
}
