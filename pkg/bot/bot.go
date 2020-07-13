package bot

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/rs/zerolog"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers"
)

const (
	// Ping event name.
	Ping = "ping"
	// IssueComment event name.
	IssueComment = "issue_comment"
)

// Config defines configuration which this bot needs to Run.
type Config struct {
	Port          string
	Hostname      string
	ServerKeyPath string
	ServerCrtPath string
	AutoTLS       bool
	CacheDir      string
	HookSecret    string
}

// Dependencies defines the dependencies of this server.
type Dependencies struct {
	Logger  zerolog.Logger
	Starter providers.Starter
}

// GaiaBot is the bot's main handler.
type GaiaBot struct {
	Dependencies
	Config
}

// Bot represents the capabilities of a Gaia bot.
type Bot interface {
	Hook(ctx context.Context) echo.HandlerFunc
}

// NewBot creates a new bot to listen for PR created events.
func NewBot(cfg Config, deps Dependencies) *GaiaBot {
	return &GaiaBot{
		Config:       cfg,
		Dependencies: deps,
	}
}

// Hook represent a github based webhook context.
type Hook struct {
	Signature string
	Event     string
	ID        string
	Payload   []byte
}

// Sender is the sender of the comment
type Sender struct {
	Login string `json:"login"`
}

// Comment is information about the comment including the URL which has to be GET
// in order to retrieve the actual comment.
type Comment struct {
	URL string `json:"url"`
}

// Issue is information about the context of the comment. It should be checked if it's a PR.
type Issue struct {
	PullRequest PullRequest `json:"pull_request,omitempty"`
	Sender      Sender      `json:"sender"`
	Comment     Comment     `json:"comment"`
}

// PullRequest is the pull request context of an issue comment if the comment happened on a PR.
type PullRequest struct {
	URL string `json:"url,omitempty"`
}

// Payload contains information about the event like, user, commit id and so on.
type Payload struct {
	Action string `json:"action"`
	Issue  Issue  `json:"issue"`
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	_, _ = computed.Write(body)
	return computed.Sum(nil)
}

func verifySignature(secret []byte, signature string, body []byte) bool {
	signaturePrefix := "sha1="
	signatureLength := 45

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	actual := make([]byte, 20)
	_, _ = hex.Decode(actual, []byte(signature[5:]))
	expected := signBody(secret, body)
	return hmac.Equal(expected, actual)
}

func parseRequest(secret []byte, req *http.Request) (Hook, error) {
	h := Hook{}

	if h.Signature = req.Header.Get("x-hub-signature"); len(h.Signature) == 0 {
		return Hook{}, errors.New("no signature")
	}

	if h.Event = req.Header.Get("x-github-event"); len(h.Event) == 0 {
		return Hook{}, errors.New("no event")
	}

	if h.Event != IssueComment {
		if h.Event == Ping {
			return Hook{Event: "ping"}, nil
		}
		return Hook{}, errors.New("invalid event")
	}

	if h.ID = req.Header.Get("x-github-delivery"); len(h.ID) == 0 {
		return Hook{}, errors.New("no event id")
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return Hook{}, err
	}

	if !verifySignature(secret, h.Signature, body) {
		return Hook{}, errors.New("invalid signature")
	}

	h.Payload = body
	return h, err
}

// GitWebHook creates a hook handler with an injected labeler.
func (b *GaiaBot) Hook(ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		h, err := parseRequest([]byte(b.Config.HookSecret), c.Request())
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		if h.Event == Ping {
			return c.NoContent(http.StatusOK)
		}

		c.Request().Header.Set("Content-type", "application/json")

		p := new(Payload)
		if err := json.Unmarshal(h.Payload, p); err != nil {
			return c.String(http.StatusBadRequest, "error in unmarshalling json payload")
		}

		if p.Action != "created" {
			return c.String(http.StatusOK, "skipped; status was not created but: "+p.Action)
		}

		if p.Issue.PullRequest.URL == "" {
			return c.String(http.StatusOK, "skipped; not in a pull request context")
		}

		if err := b.Dependencies.Starter.Start(ctx, p.Issue.Sender.Login, p.Issue.Comment.URL, p.Issue.PullRequest.URL); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, "successfully processed event")
	}
}
