package providers

import (
	"context"
)

// Commenter will add a comment of failed, success, ack to a PR.
type Commenter interface {
	AddComment(ctx context.Context, number string, comment string) error
}
