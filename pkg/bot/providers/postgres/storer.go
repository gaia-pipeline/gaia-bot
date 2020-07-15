package postgres

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/models"
)

// Config has the configuration options for the starter
type Config struct {
	Hostname string
	Database string
	Username string
	Password string
}

// Dependencies defines the dependencies of this bot
type Dependencies struct {
	Logger zerolog.Logger
}

// Store is a postgres based store.
type Store struct {
	Config
	Dependencies
}

// NewPostgresStore creates a new PostgresStore
func NewPostgresStore(cfg Config, deps Dependencies) *Store {
	return &Store{Config: cfg, Dependencies: deps}
}

// GetUser returns the full user entity.
func (s *Store) GetUser(ctx context.Context, handle string) (*models.User, error) {
	return &models.User{
		Handle:   "Skarlso",
		Commands: []string{"test"},
	}, nil
}
