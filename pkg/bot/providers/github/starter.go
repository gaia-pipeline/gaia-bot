package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
	Store  providers.Storer
	Logger zerolog.Logger
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
	log.Info().Msg("Starting update...")
	// user has access and we know the user... launch the update.
	return nil
}

// extractComment gets a comment from a comment url.
func (s *Starter) extractComment(ctx context.Context, commentUrl string) (string, error) {
	s.Logger.Debug().Str("comment-url", commentUrl).Msg("extacting comment from url")
	resp, err := http.Get(commentUrl)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to get comment body")
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error().Err(err).Msg("Failed to get response body")
		return "", err
	}
	c := comment{}
	if err := json.Unmarshal(body, &c); err != nil {
		s.Logger.Error().Err(err).Msg("Failed to unmarshall comment body")
		return "", err
	}
	return c.Body, nil
}
