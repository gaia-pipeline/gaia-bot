package providers

import "context"

// Executioner defines an execution method.
type Executioner interface {
	// Run a command over an executioning medium.
	Execute(ctx context.Context, script string, args map[string]string) error
}
