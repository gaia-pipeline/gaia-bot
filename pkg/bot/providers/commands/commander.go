package commands

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/auth"
)

const (
	sumFormat       = "%s-%d-%d"
	imageTag        = "gaiapipeline/testing:%x"
	fetchPrFilenane = "/usr/local/bin/fetch_pr.sh"
	pushTagFilename = "/usr/local/bin/push_tag.sh"
)

const (
	helpMessage = "Hello, I'm Gaia Bot. I help the team build, test and manage PRs and other issues.\n" +
		"These are my available commands:\n" +
		"`/test`: Apply the current PR to our testing infrastructure at https://gaia.cronohub.org.\n" +
		"Example: `/test`"
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
	Converter   providers.EnvironmentConverter
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
func (c *Commander) Test(ctx context.Context, owner string, repoURL string, repo string, number int, branch string, commentID int) {
	log := c.Logger.With().Int("number", number).Str("branch", branch).Logger()
	if err := c.Commenter.AddComment(ctx, owner, repo, number, "Command received. Building new test image."); err != nil {
		log.Error().Err(err).Msg("Failed to add ack comment.")
		return
	}

	// Create a tag based on the branch, number and commentID
	tagSum := fmt.Sprintf(sumFormat, branch, number, commentID)
	sum := sha256.Sum256([]byte(tagSum))
	tag := fmt.Sprintf(imageTag, sum)

	// Run the docker build over SSH
	script, err := ioutil.ReadFile(fetchPrFilenane)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read script.")
		return
	}
	n := strconv.Itoa(number)
	dockerToken, err := c.Dependencies.Converter.LoadValueFromFile(c.Auth.DockerToken)
	if err != nil {
		log.Error().Err(err).Msg("Failed LoadValueFromFile.")
		return
	}
	dockerUsername, err := c.Dependencies.Converter.LoadValueFromFile(c.Auth.DockerUsername)
	if err != nil {
		log.Error().Err(err).Msg("Failed LoadValueFromFile.")
		return
	}
	if err := c.Executioner.Execute(ctx, string(script), map[string]string{
		"<repo_url_replace>":        repoURL,
		"<tag_replace>":             tag,
		"<pr_replace>":              n,
		"<branch_replace>":          branch,
		"<docker_token_replace>":    dockerToken,
		"<docker_username_replace>": dockerUsername,
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
	githubToken, err := c.Dependencies.Converter.LoadValueFromFile(c.Auth.GithubToken)
	if err != nil {
		log.Error().Err(err).Msg("Failed LoadValueFromFile.")
		return
	}
	githubUsername, err := c.Dependencies.Converter.LoadValueFromFile(c.Auth.GithubUsername)
	if err != nil {
		log.Error().Err(err).Msg("Failed LoadValueFromFile.")
		return
	}
	if err := c.Executioner.Execute(ctx, string(script), map[string]string{
		"<repo_replace>":         c.InfraRepo,
		"<tag_replace>":          tag,
		"<git_token_replace>":    githubToken,
		"<git_username_replace>": githubUsername,
	}); err != nil {
		log.Error().Err(err).Msg("Failed to execute command.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to push new code to gaia infrastructure repository."); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}
	if err := c.Commenter.AddComment(ctx, owner, repo, number, "New version pushed. Please wait for 5 minutes for it to propagate."); err != nil {
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
