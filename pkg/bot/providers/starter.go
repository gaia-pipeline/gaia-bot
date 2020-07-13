package providers

import (
	"context"
)

// Starter responds to comments made on a PR. If the person is authorized to start a test
// it will start deploying the PR's code as a new docker image and put it under the
// deployed test environment.
type Starter interface {
	// Start to deploy a new test environment
	Start(ctx context.Context, handle string, commentURL string, pullURL string) error
}
