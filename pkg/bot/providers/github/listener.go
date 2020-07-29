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

// Config has the configuration options for the listener
type Config struct {
}

// Dependencies defines the dependencies of this bot
type Dependencies struct {
	Store     providers.Storer
	Commander providers.Commander
	Logger    zerolog.Logger
}

// Listener is a github based bot.
type Listener struct {
	Config
	Dependencies
}

// NewGithubListener creates a new GithubListener
func NewGithubListener(cfg Config, deps Dependencies) *Listener {
	return &Listener{Config: cfg, Dependencies: deps}
}

// Listen starts to listen on PR events and respond accordingly.
func (s *Listener) Listen(ctx context.Context, handle string, commentURL string, pullURL string) error {
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

	repo, err := s.extractPullRequestInfo(ctx, pullURL)
	if err != nil {
		log.Error().Err(err).Msg("Failed to extract pull request information.")
		return err
	}

	// launch the Command with a new context and return
	switch cmd {
	case "test":
		log.Info().Msg("Starting update...")
		go s.Dependencies.Commander.Test(context.Background(), repo.Base.Repo.Owner.Login, repo.Base.Repo.URL, repo.Base.Repo.Name, repo.Number, repo.Head.Ref)
	default:
		return fmt.Errorf("command %s not found", cmd)
	}
	return nil
}

// head is the head of the git repository
type head struct {
	Ref string `json:"ref"`
}

// owner has the owner information of the forked repository
type owner struct {
	Login string `json:"login"`
}

// repoInfo is the repo part of the base of the forked repository
type repoInfo struct {
	Owner owner  `json:"owner"`
	Name  string `json:"name"`
	URL   string `json:"html_url"`
}

// base is the base of the forked repository
type base struct {
	Repo repoInfo `json:"repo"`
}

// repo is a git repository
type repo struct {
	Head   head `json:"head"`
	Number int  `json:"number"`
	Base   base `json:"base"`
}

// extractPullRequestInfo return with the details of the repository.
func (s *Listener) extractPullRequestInfo(ctx context.Context, pullUrl string) (repo, error) {
	r := repo{}
	if err := s.get(ctx, pullUrl, &r); err != nil {
		return repo{}, err
	}
	return r, nil
}

// extractComment gets a comment from a comment url.
func (s *Listener) extractComment(ctx context.Context, commentUrl string) (string, error) {
	s.Logger.Debug().Str("comment-url", commentUrl).Msg("extacting comment from url")
	c := comment{}
	if err := s.get(ctx, commentUrl, &c); err != nil {
		return "", err
	}
	return c.Body, nil
}

func (s *Listener) get(ctx context.Context, url string, v interface{}) error {
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
