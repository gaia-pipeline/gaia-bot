// commands will handle all the commands this bot supports.
package commands

import (
	"context"
	"io/ioutil"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
)

// Config has the configuration options for the commander
type Config struct {
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
func (c *Commander) Test(ctx context.Context, pullRequest string) {
	log := c.Logger.With().Str("pull-request", pullRequest).Logger()
	tmp, err := ioutil.TempDir("checkout", "gaia")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create temp directory to checkout pr.")
		if err := c.Commenter.AddComment(ctx, "", "Failed"); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}
	log.Debug().Str("tmp", tmp).Msg("Temp dir created.")
	// pull gaia main repo
	// do a fetch
	// switch to branch
}
