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

const (
	imageTag        = "gaiapipeline/testing:%s-%d"
	fetchPrFilenane = "/usr/local/bin/fetch_pr.sh"
	pushTagFilename = "/usr/local/bin/push_tag.sh"
)

const (
	helpMessage = "Hello, I'm Gaia Bot. I help the team build, test and manage PRs and other issues.\n" +
		"These are my available commands:\n" +
		"`/test`: Apply the current PR to our testing infrastructure at https://gaia.cronohub.org.\n" +
		"Example: `/test`, `/test gaiapipeline/testing:tag-name`"
)

// Config has the configuration options for the commander
type Config struct {
	Auth      auth.Config
	InfraRepo string
}

// Dependencies defines the dependencies of this command
type Dependencies struct {
	Logger      zerolog.Logger
	Commenter   providers.Commenter
	Executioner providers.Executer
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
func (c *Commander) Test(ctx context.Context, owner string, repoURL string, repo string, number int, branch string, paramTag string) {
	log := c.Logger.With().Int("number", number).Str("branch", branch).Logger()
	if err := c.Commenter.AddComment(ctx, owner, repo, number, "Command received. Building new test image."); err != nil {
		log.Error().Err(err).Msg("Failed to add ack comment.")
		return
	}
	tag := fmt.Sprintf(imageTag, branch, number)
	if paramTag != "" {
		tag = paramTag
	}

	// Run the docker build over SSH
	script, err := ioutil.ReadFile(fetchPrFilenane)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read script.")
		return
	}
	n := strconv.Itoa(number)
	if err := c.Executioner.Execute(ctx, string(script), map[string]string{
		"<repo_url_replace>":        repoURL,
		"<tag_replace>":             tag,
		"<pr_replace>":              n,
		"<branch_replace>":          branch,
		"<docker_token_replace>":    c.Auth.DockerToken,
		"<docker_username_replace>": c.Auth.DockerUsername,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to fetch and build pr.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to fetch and build pr."); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}

	// Run the pusher
	script, err = ioutil.ReadFile(pushTagFilename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read script.")
		return
	}
	if err := c.Executioner.Execute(ctx, string(script), map[string]string{
		"<repo_replace>":         c.InfraRepo,
		"<tag_replace>":          tag,
		"<git_token_replace>":    c.Auth.GithubToken,
		"<git_username_replace>": c.Auth.GithubUsername,
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

// Help displays information about what commands are available and their usage.
func (c *Commander) Help(ctx context.Context, owner string, repo string, number int) {
	if err := c.Commenter.AddComment(ctx, owner, repo, number, helpMessage); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to add comment.")
		return
	}
}
