// commands will handle all the commands this bot supports.
package commands

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
)

// Config has the configuration options for the commander
type Config struct {
	Token    string
	Username string
}

// Dependencies defines the dependencies of this command
type Dependencies struct {
	Logger    zerolog.Logger
	Commenter providers.Commenter
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
func (c *Commander) Test(ctx context.Context, owner string, repo string, number int, branch string) {
	log := c.Logger.With().Int("number", number).Str("branch", branch).Logger()
	if err := c.Commenter.AddComment(ctx, owner, repo, number, "Command received. Running Test."); err != nil {
		log.Error().Err(err).Msg("Failed to add ack comment.")
		return
	}
	// Just commit changes to the flux repository with the new image tag.
}
