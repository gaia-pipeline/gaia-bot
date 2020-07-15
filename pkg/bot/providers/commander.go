package providers

import (
	"context"
)

// Commander represents the commands which this bot is capable of performing
type Commander interface {
	HandleCommand(ctx context.Context, cmd string) error
}
