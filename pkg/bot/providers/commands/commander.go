// commands will handle all the commands this bot supports.
package commands

import (
	"context"

	"github.com/rs/zerolog"
)

// command is a command that a commander can perform
type command func(ctx context.Context) error

// CommandMap represents a mapping between command names and the command that can be performed.
var commandMap = map[string]command{}

// Config has the configuration options for the starter
type Config struct {
}

// Dependencies defines the dependencies of this bot
type Dependencies struct {
	Logger zerolog.Logger
}

// Commander is a bot commander.
type Commander struct {
	Config
	Dependencies
}

// NewCommander creates a new GithubStarter
func NewCommander(cfg Config, deps Dependencies) *Commander {
	c := &Commander{Config: cfg, Dependencies: deps}
	commandMap["test"] = c.Test
	return c
}

// HandleCommand takes a command and performs the appropriate action.
func (c *Commander) HandleCommand(ctx context.Context, cmd string) error {
	return nil
}

// Test deploys a PR to a test environment
func (c *Commander) Test(ctx context.Context) error {
	return nil
}
