package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
)

const commentPrefix = "/"

// user is a user inside the comment.
type user struct {
	Login string `json:"login"`
}

// comment represents a comment json structure returned by github.
type comment struct {
	User user   `json:"user"`
	Body string `json:"body"`
}

// Config has the configuration options for the starter
type Config struct {
}

// Dependencies defines the dependencies of this bot
type Dependencies struct {
	Store     providers.Storer
	Commander providers.Commander
	Logger    zerolog.Logger
}

// Starter is a github based bot.
type Starter struct {
	Config
	Dependencies
}

// NewGithubStarter creates a new GithubStarter
func NewGithubStarter(cfg Config, deps Dependencies) *Starter {
	return &Starter{Config: cfg, Dependencies: deps}
}

// Start starts an update if the user has access.
func (s *Starter) Start(ctx context.Context, handle string, commentURL string, pullURL string) error {
	log := s.Logger.With().Str("handle", handle).Logger()
	comment, err := s.extractComment(ctx, commentURL)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(comment, commentPrefix) {
		log.Debug().Str("comment", comment).Msg("Ignored comment as it doesn't have the bot prefix.")
		return nil
	}
	cmd := strings.TrimPrefix(comment, commentPrefix)
	log.Debug().Msg("Checking if handle has access to the command...")
	user, err := s.Store.GetUser(ctx, handle)
	if err != nil {
		log.Error().Err(err).Msg("GetUser failed")
		return err
	}
	if user == nil {
		log.Error().Msg("User not found")
		return errors.New("user not found")
	}
	hasCommand := false
	for _, c := range user.Commands {
		if c == cmd {
			hasCommand = true
			break
		}
	}
	if !hasCommand {
		return fmt.Errorf("user doesn't have access to command %s", cmd)
	}

	id, ref, err := s.extractPullRequestInfo(ctx, pullURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to extract pull request information.")
		return err
	}

	// launch the Command with a new context and return
	switch cmd {
	case "test":
		log.Info().Msg("Starting update...")
		go s.Dependencies.Commander.Test(context.Background(), id, ref)
	default:
		return fmt.Errorf("command %s not found", cmd)
	}
	return nil
}

// head is the head of the git repository
type head struct {
	Ref string `json:"ref"`
}

// repo is a git repository
type repo struct {
	Head head `json:"head"`
}

// extractPullRequestInfo return with the details of the repository.
func (s *Starter) extractPullRequestInfo(ctx context.Context, pullUrl string) (string, string, error) {
	split := strings.Split(pullUrl, "/")
	if len(split) < 1 {
		s.Logger.Error().Str("url", pullUrl).Msg("Split length below 1.")
		return "", "", fmt.Errorf("split length below 1. was %d", len(split))
	}
	id := split[len(split)-1]
	r := repo{}
	if err := s.get(ctx, pullUrl, &r); err != nil {
		return "", "", err
	}
	return id, r.Head.Ref, nil
}

// extractComment gets a comment from a comment url.
func (s *Starter) extractComment(ctx context.Context, commentUrl string) (string, error) {
	s.Logger.Debug().Str("comment-url", commentUrl).Msg("extacting comment from url")
	c := comment{}
	if err := s.get(ctx, commentUrl, &c); err != nil {
		return "", err
	}
	return c.Body, nil
}

func (s *Starter) get(ctx context.Context, url string, v interface{}) error {
	log := s.Logger.With().Str("url", url).Logger()
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get body")
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get response body")
		return err
	}
	if err := json.Unmarshal(body, &v); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshall body")
		return err
	}
	return nil
}
