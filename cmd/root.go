package main

import (
	"context"
	"flag"
	"os"

	"github.com/gaia-pipeline/gaia-bot/pkg/bot/providers/auth"

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
		commander commands.Config
		auth      auth.Config
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
	flag.StringVar(&rootArgs.store.Hostname, "gaia-bot-db-hostname", "localhost", "--gaia-bot-db-hostname localhost")
	flag.StringVar(&rootArgs.store.Database, "gaia-bot-db-database", "bot", "--gaia-bot-db-database bot")
	flag.StringVar(&rootArgs.store.Username, "gaia-bot-db-username", "bot", "--gaia-bot-db-username bot")
	flag.StringVar(&rootArgs.store.Password, "gaia-bot-db-password", "password123", "--gaia-bot-db-password password123")
	flag.StringVar(&rootArgs.bot.HookSecret, "hook-secret", "", "--hook-secret asdf")
	flag.StringVar(&rootArgs.auth.GithubToken, "githu-token", "", "--github-token asdf")
	flag.StringVar(&rootArgs.auth.GithubUsername, "githu-username", "", "--github-username asdf")
	flag.StringVar(&rootArgs.auth.DockerToken, "docker-token", "", "--docker-token asdf")
	flag.StringVar(&rootArgs.auth.DockerUsername, "docker-username", "", "--docker-username asdf")
	flag.StringVar(&rootArgs.commander.InfraRepo, "infra-repository", "github.com/gaia-pipeline/gaia-infrastructure.git", "--infra-repository github.com/")
	flag.Parse()
}

func main() {
	log := zerolog.New(os.Stderr)
	storer := postgres.NewPostgresStore(rootArgs.store, postgres.Dependencies{
		Logger: log,
	})
	rootArgs.commenter.Auth = rootArgs.auth
	commenter := commenter.NewCommenter(rootArgs.commenter, commenter.Dependencies{Logger: log})
	rootArgs.commander.Auth = rootArgs.auth
	commander := commands.NewCommander(rootArgs.commander, commands.Dependencies{
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
