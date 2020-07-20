// commands will handle all the commands this bot supports.
package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
)

const imageTag = "gaiapipeline/testing:%s-%d"

// Config has the configuration options for the commander
type Config struct {
	GitToken       string
	GitUsername    string
	DockerToken    string
	DockerUsername string
	InfraRepo      string
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
	tmp, err := ioutil.TempDir("checkout", "build")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create temp directory to checkout pr.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to checkout the PR"); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}
	tag := fmt.Sprintf(imageTag, branch, number)

	// Run the docker build
	n := strconv.Itoa(number)
	cmd := exec.Command("/usr/local/bin/fetch_pr.sh")
	cmd.Env = append(os.Environ(), ""+
		"TAG="+tag,
		"PR_NUMBER="+n,
		"REPOSITORY="+repo,
		"BRANCH="+branch,
		"FOLDER="+tmp,
		"DOCKER_TOKEN="+c.DockerToken,
		"DOCKER_USERNAME="+c.DockerUsername,
	)
	if err := cmd.Run(); err != nil {
		log.Error().Err(err).Msg("Failed to run fetcher.")
		if err := c.Commenter.AddComment(ctx, owner, repo, number, "Failed to fetch PR."); err != nil {
			log.Error().Err(err).Msg("Failed to add comment.")
			return
		}
		return
	}

	// Run the pusher
	log.Debug().Msgf("pusing image tag %s", tag)

	cmd = exec.Command("/usr/local/bin/push_tag.sh")
	cmd.Env = append(os.Environ(), ""+
		"REPO="+c.InfraRepo,
		"TAG="+tag,
		"FOLDER="+tmp,
		"GIT_TOKEN="+c.GitToken,
		"GIT_USERNAME="+c.GitUsername,
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
