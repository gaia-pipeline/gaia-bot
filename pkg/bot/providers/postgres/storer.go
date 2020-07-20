package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/models"
)

const timeoutForTransactions = 1 * time.Minute

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
	conn, err := s.connect()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, timeoutForTransactions)
	defer cancel()

	defer conn.Close(ctx)
	var (
		storedHandle string
		commands     string
	)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "select handle, commands from users where handle = $1", handle).Scan(&storedHandle, &commands)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	splitCommands := strings.Split(commands, ",")
	return &models.User{
		Handle:   storedHandle,
		Commands: splitCommands,
	}, nil
}

func (s Store) connect() (*pgx.Conn, error) {
	url := fmt.Sprintf("postgresql://%s/%s?user=%s&password=%s", s.Hostname, s.Database, s.Username, s.Password)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to connect to the database")
		return nil, err
	}
	return conn, nil
}
