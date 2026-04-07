package service

import (
	"context"
	"fmt"

	"github.com/stkrfx/dg-backend/internal/model"
	"github.com/stkrfx/dg-backend/internal/repository"
)

type ClientService struct {
	repo *repository.ClientRepository
}

func NewClientService(r *repository.ClientRepository) *ClientService {
	return &ClientService{repo: r}
}

// RegisterClient now takes a context and returns errors properly
func (s *ClientService) RegisterClient(ctx context.Context, email string, name string) error {
	// Professional Tip: In a real app, you would add validation here
	if email == "" {
		return fmt.Errorf("email is required")
	}

	// Create the client model
	// We no longer calculate ID manually (Supabase/Postgres does this for us)
	client := &model.Client{
		Email: email,
		Name:  name,
	}

	// Pass the context and the pointer to the repository
	err := s.repo.Create(ctx, client)
	if err != nil {
		// Wrap the error to provide context on where it failed
		return fmt.Errorf("service.RegisterClient: %w", err)
	}

	return nil
}

// FetchAllClients allows the handler to get the full list
func (s *ClientService) FetchAllClients(ctx context.Context) ([]model.Client, error) {
	clients, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("service.FetchAllClients: %w", err)
	}
	return clients, nil
}
