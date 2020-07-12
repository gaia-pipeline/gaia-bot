package server

import (
	"context"
)

// Config is the configuration of the bot's main cycle.
type Config struct {
}

// GaiaBot is a server.
type GaiaBot struct {
	Config
}

// Server defines a server which runs and accepts requests and monitors
// Gaia PRs.
type Server interface {
	Run(context.Context) error
}

// Run starts up the listening for PR actions.
func (g *GaiaBot) Run(ctx context.Context) error {
	return nil
}
