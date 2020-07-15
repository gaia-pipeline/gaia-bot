package providers

import (
	"context"

	"github.com/gaia-pipeline/gaia-bot/pkg/models"
)

// Storer represents the functionality of a storage medium which deals with
// retrieving and managing users who have access to the gaia-bot.
type Storer interface {
	// GetUser returns the full user entity.
	GetUser(ctx context.Context, handle string) (*models.User, error)
}
