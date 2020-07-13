package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot"
)

// Config is the configuration of the bot's main cycle.
type Config struct {
	Port          string
	Hostname      string
	ServerKeyPath string
	ServerCrtPath string
	AutoTLS       bool
	CacheDir      string
	HookSecret    string
}

// GaiaBot is a server.
type GaiaBot struct {
	Config
	Logger zerolog.Logger
	Bot    bot.Bot
}

// Server defines a server which runs and accepts requests and monitors
// Gaia PRs.
type Server interface {
	Run(context.Context) error
}

// Run starts up the listening for PR actions.
func (s *GaiaBot) Run(ctx context.Context) error {
	s.Logger.Info().Msg("Start listening...")
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/hook", s.Bot.Hook(ctx))

	hostPort := fmt.Sprintf("%s:%s", s.Config.Hostname, s.Config.Port)

	// Start TLS with certificate paths
	if len(s.Config.ServerKeyPath) > 0 && len(s.Config.ServerCrtPath) > 0 {
		e.Pre(middleware.HTTPSRedirect())
		return e.StartTLS(hostPort, s.Config.ServerCrtPath, s.Config.ServerKeyPath)
	}

	// Start Auto TLS server
	if s.Config.AutoTLS {
		if len(s.Config.CacheDir) < 1 {
			return errors.New("cache dir must be provided if autoTLS is enabled")
		}
		e.Pre(middleware.HTTPSRedirect())
		e.AutoTLSManager.Cache = autocert.DirCache(s.Config.CacheDir)
		return e.StartAutoTLS(hostPort)
	}
	// Start regular server
	return e.Start(hostPort)
}
