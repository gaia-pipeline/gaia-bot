package providers

import (
	"context"
)

// Commander represents the commands which this bot is capable of performing
type Commander interface {
	Test(ctx context.Context, pullRequest string) error
}
