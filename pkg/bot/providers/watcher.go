package providers

import (
	"context"
)

// Watcher watches a pull-request on github.
type Watcher interface {
	// Start watching pull requests on Gaia's github page.
	Watch(ctx context.Context) error
}
