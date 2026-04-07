package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/stkrfx/dg-backend/internal/model"
)

type ClientRepository struct {
	db *sql.DB
}

// NewClientRepository now accepts the shared DB connection
func NewClientRepository(db *sql.DB) *ClientRepository {
	return &ClientRepository{db: db}
}

// Create inserts a client and updates the struct with the new ID from Supabase
func (r *ClientRepository) Create(ctx context.Context, c *model.Client) error {
	query := `
		INSERT INTO clients (name, email, created_at) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	// Use QueryRowContext to get the auto-generated ID back
	err := r.db.QueryRowContext(ctx, query, c.Name, c.Email, time.Now()).Scan(&c.ID)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	return nil
}

// GetAll fetches all clients from the Supabase table
func (r *ClientRepository) GetAll(ctx context.Context) ([]model.Client, error) {
	query := `SELECT id, name, email, created_at FROM clients`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clients: %w", err)
	}
	defer rows.Close() // Essential: returns the connection to the pool

	var clients []model.Client
	for rows.Next() {
		var c model.Client
		// Scan fields in the exact order of the SELECT statement
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning client row: %w", err)
		}
		clients = append(clients, c)
	}

	// Check for any errors that occurred during the loop
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}
