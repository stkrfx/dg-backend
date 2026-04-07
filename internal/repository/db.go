package repository

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var fs embed.FS

func NewDB(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	// Professional connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Verify connection immediately
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Run automatic migrations
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Println("Successfully connected to Supabase Postgres")
	return db, nil
}

// runMigrations automatically applies any pending SQL files in the migrations folder
func runMigrations(db *sql.DB) error {
	// 1. Read the embedded SQL files
	sourceDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	// 2. Connect the migration engine to our Postgres DB
	dbDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// 3. Initialize the migrate instance
	m, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return err
	}

	// 4. Apply all 'Up' migrations
	err = m.Up()

	// 'ErrNoChange' just means the DB is already up to date, which is completely fine
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("Database schemas are verified and up to date.")
	return nil
}
