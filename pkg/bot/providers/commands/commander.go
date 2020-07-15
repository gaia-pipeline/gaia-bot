// commands will handle all the commands this bot supports.
package commands

import (
	"context"

	"github.com/rs/zerolog"
)

// Config has the configuration options for the commander
type Config struct {
}

// Dependencies defines the dependencies of this command
type Dependencies struct {
	Logger zerolog.Logger
}

// Commander is a bot commander.
type Commander struct {
	Config
	Dependencies
}

// NewCommander creates a new Commander
func NewCommander(cfg Config, deps Dependencies) *Commander {
	c := &Commander{Config: cfg, Dependencies: deps}
	return c
}

// Test will deploy the pull request from the context of the comment made.
func (c *Commander) Test(ctx context.Context, pullRequest string) error {
	return nil
}
