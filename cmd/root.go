package main

import (
	"context"
	"flag"
	"os"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot"
	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/commands"
	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/commenter"
	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/github"
	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/postgres"
	"github.com/gaia-pipeline/gaia-bot/pkg/server"
)

var (
	rootArgs struct {
		devMode   bool
		bot       bot.Config
		store     postgres.Config
		server    server.Config
		commenter commenter.Config
		debug     bool
	}
)

func init() {
	flag.BoolVar(&rootArgs.server.AutoTLS, "auto-tls", false, "--auto-tls")
	flag.BoolVar(&rootArgs.debug, "debug", false, "--debug")
	flag.StringVar(&rootArgs.server.CacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&rootArgs.server.ServerKeyPath, "server-key-path", "", "--server-key-path /home/user/.server/server.key")
	flag.StringVar(&rootArgs.server.ServerCrtPath, "server-crt-path", "", "--server-crt-path /home/user/.server/server.crt")
	flag.StringVar(&rootArgs.server.Port, "port", "9998", "--port 443")
	flag.StringVar(&rootArgs.server.Hostname, "hostname", "", "--hostname gaia-bot.org")
	flag.StringVar(&rootArgs.store.Hostname, "staple-db-hostname", "localhost", "--gaia-bot-db-hostname localhost")
	flag.StringVar(&rootArgs.store.Database, "staple-db-database", "bots", "--gaia-bot-db-database staples")
	flag.StringVar(&rootArgs.store.Username, "staple-db-username", "bot", "--gaia-bot-db-username staple")
	flag.StringVar(&rootArgs.store.Password, "staple-db-password", "password123", "--gaia-bot-db-password password123")
	flag.StringVar(&rootArgs.bot.HookSecret, "hook-secret", "", "--hook-secret asdf")
	flag.StringVar(&rootArgs.commenter.Token, "github-token", "", "--github-token asdf")
	flag.Parse()
}

func main() {
	log := zerolog.New(os.Stderr)
	storer := postgres.NewPostgresStore(rootArgs.store, postgres.Dependencies{
		Logger: log,
	})
	commenter := commenter.NewCommenter(rootArgs.commenter, commenter.Dependencies{Logger: log})
	commander := commands.NewCommander(commands.Config{}, commands.Dependencies{
		Logger:    log,
		Commenter: commenter,
	})
	starter := github.NewGithubStarter(github.Config{}, github.Dependencies{
		Logger:    log,
		Store:     storer,
		Commander: commander,
	})
	gaiaBot := bot.NewBot(rootArgs.bot, bot.Dependencies{
		Logger:  log,
		Starter: starter,
	})

	botServer := server.NewServer(rootArgs.server, server.Dependencies{
		Logger: log,
		Bot:    gaiaBot,
	})
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error { return botServer.Run(ctx) })
	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Failed to run")
	}
}
