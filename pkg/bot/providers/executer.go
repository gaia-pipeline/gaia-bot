package providers

import "context"

// Executer defines an execution method.
type Executer interface {
	// Run a command over an executioning medium.
	Execute(ctx context.Context, script string, args map[string]string) error
}
