package github

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog"
)

// user is a user inside the comment.
type user struct {
	Login string `json:"login"`
}

// comment represents a comment json structure returned by github.
type comment struct {
	User user   `json:"user"`
	Body string `json:"body"`
}

// extractComment gets a comment from a comment url.
func extractComment(ctx context.Context, log zerolog.Logger, commentUrl string) (string, error) {
	resp, err := http.Get("http://example.com/")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get comment body")
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get response body")
		return "", err
	}
	c := comment{}
	if err := json.Unmarshal(body, &c); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshall comment body")
		return "", err
	}
	return c.Body, nil
}
