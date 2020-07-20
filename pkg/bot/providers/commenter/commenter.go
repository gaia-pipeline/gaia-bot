package commenter

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/auth"
)

// Config has the configuration options for the commenter
type Config struct {
	Auth auth.Config
}

// Dependencies defines the dependencies of this commenter
type Dependencies struct {
	Logger zerolog.Logger
}

// PullRequestsService is an interface defining the Wrapper Interface
// needed to test the github client.
type PullRequestsService interface {
	CreateComment(ctx context.Context, owner string, repo string, number int, comment *github.IssueComment) (*github.IssueComment, *github.Response, error)
}

// Commenter is a github based commenter.
type Commenter struct {
	Config
	Dependencies
	Client *github.IssuesService
}

// NewCommenter creates a new Commenter
func NewCommenter(cfg Config, deps Dependencies) *Commenter {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.Auth.GithubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return &Commenter{Config: cfg, Dependencies: deps, Client: client.Issues}
}

// AddComment adds a comment to a PR to show progress of the bot.
func (c *Commenter) AddComment(ctx context.Context, owner string, repo string, number int, comment string) error {
	if _, _, err := c.Client.CreateComment(ctx, owner, repo, number, &github.IssueComment{Body: &comment}); err != nil {
		c.Logger.Error().Err(err).Msg("Failed to add comment")
		return err
	}
	return nil
}
