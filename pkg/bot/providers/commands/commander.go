// commands will handle all the commands this bot supports.
package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
)

const imageTag = "gaiapipeline/testing:%s-%d"

// Config has the configuration options for the commander
type Config struct {
	Token     string
	Username  string
	InfraRepo string
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
	tmp, err := ioutil.TempDir("checkout", "gaia-infrastructure")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create temp directory to checkout pr.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to checkout the PR"); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}
	// Just commit changes to the flux repository with the new image tag.
	tag := fmt.Sprintf(imageTag, branch, number)
	log.Debug().Msgf("pusing image tag %s", tag)

	cmd := exec.Command("/usr/local/bin/push_tag.sh")
	cmd.Env = append(os.Environ(), ""+
		"REPO="+c.InfraRepo,
		"TAG="+tag,
		"FOLDER="+tmp,
	)
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run pusher.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to push new code to gaia infrastructure repository."); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}
	if err := c.Commenter.AddComment(ctx, owner, repo, number, "The new test version has been pushed."); err != nil {
		log.Error().Err(err).Msg("Failed to add comment.")
		return
	}
}
