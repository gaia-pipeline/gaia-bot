package sshexecutioner

import (
	"context"

	"github.com/rs/zerolog"
)

// Config has the configuration options for the commenter
type Config struct {
	Address  string
	Username string
}

// Dependencies defines the dependencies of this commenter
type Dependencies struct {
	Logger zerolog.Logger
}

// SSHExec is an executioner which can execute commands over SSH.
// This assumes that the current runtime environment has SSH access to the machine.
type SSHExec struct {
	Config
	Dependencies
}

// NewSshExecutioner creates a new SSHExec
func NewSshExecutioner(cfg Config, deps Dependencies) *SSHExec {
	return nil
}

// Executioner defines an execution method.
func (e *SSHExec) Execute(ctx context.Context, script string, args map[string]string) error {
	return nil
}
