package providers

import (
	"context"
)

// Processor processes a comment event created on a PR.
type Processor interface {
	// Process processes a single pr comment.
	Process(ctx context.Context, handle string, commentURL string, pullURL string) error
}
