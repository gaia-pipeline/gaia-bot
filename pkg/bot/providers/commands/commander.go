// commands will handle all the commands this bot supports.
package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/auth"
)

const imageTag = "gaiapipeline/testing:%s-%d"

// Config has the configuration options for the commander
type Config struct {
	Auth      auth.Config
	InfraRepo string
}

// Dependencies defines the dependencies of this command
type Dependencies struct {
	Logger      zerolog.Logger
	Commenter   providers.Commenter
	Executioner providers.Executioner
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
	if err := c.Commenter.AddComment(ctx, owner, repo, number, "Command received. Building new test image."); err != nil {
		log.Error().Err(err).Msg("Failed to add ack comment.")
		return
	}
	tmp, err := ioutil.TempDir("", "build")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create temp directory to checkout pr.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to checkout the PR"); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}
	tag := fmt.Sprintf(imageTag, branch, number)

	// Run the docker build over SSH
	n := strconv.Itoa(number)
	if err := c.Executioner.Execute(ctx, "/usr/local/bin/fetch_pr.sh", map[string]string{
		"TAG":             tag,
		"PR_NUMBER":       n,
		"REPOSITORY":      repo,
		"BRANCH":          branch,
		"FOLDER":          tmp,
		"DOCKER_TOKEN":    c.Auth.DockerToken,
		"DOCKER_USERNAME": c.Auth.DockerUsername,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to fetch and build pr.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to fetch and build pr."); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}

	// Run the pusher
	if err := c.Executioner.Execute(ctx, "/usr/local/bin/push_tag.sh", map[string]string{
		"REPO":         c.InfraRepo,
		"TAG":          tag,
		"FOLDER":       tmp,
		"GIT_TOKEN":    c.Auth.GithubToken,
		"GIT_USERNAME": c.Auth.GithubUsername,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to execute command.")
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
