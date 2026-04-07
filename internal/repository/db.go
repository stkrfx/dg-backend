package repository

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

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

	log.Println("Successfully connected to Supabase Postgres")
	return db, nil
}
