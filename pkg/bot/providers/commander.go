package providers

import (
	"context"
)

// Commander represents the commands which this bot is capable of performing
type Commander interface {
	Test(ctx context.Context, owner string, repoURL string, repo string, number int, branch string, commentID int)
	Help(ctx context.Context, owner string, repo string, number int)
}
