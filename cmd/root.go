package main

import (
	"flag"
)

var (
	rootArgs struct {
		devMode       bool
		autoTLS       bool
		cacheDir      string
		serverKeyPath string
		serverCrtPath string
		port          string
		hostname      string
		hookSecret    string
		database      struct {
			hostname string
			database string
			username string
			password string
		}
		debug bool
	}
)

func init() {
	flag.BoolVar(&rootArgs.devMode, "dev", false, "--dev")
	flag.BoolVar(&rootArgs.autoTLS, "auto-tls", false, "--auto-tls")
	flag.BoolVar(&rootArgs.debug, "debug", false, "--debug")
	flag.StringVar(&rootArgs.cacheDir, "cache-dir", "", "--cache-dir /home/user/.server/.cache")
	flag.StringVar(&rootArgs.serverKeyPath, "server-key-path", "", "--server-key-path /home/user/.server/server.key")
	flag.StringVar(&rootArgs.serverCrtPath, "server-crt-path", "", "--server-crt-path /home/user/.server/server.crt")
	flag.StringVar(&rootArgs.port, "port", "9998", "--port 443")
	flag.StringVar(&rootArgs.hostname, "hostname", "", "--hostname gaia-bot.org")
	flag.StringVar(&rootArgs.database.hostname, "staple-db-hostname", "localhost", "--gaia-bot-db-hostname localhost")
	flag.StringVar(&rootArgs.database.database, "staple-db-database", "bots", "--gaia-bot-db-database staples")
	flag.StringVar(&rootArgs.database.username, "staple-db-username", "bot", "--gaia-bot-db-username staple")
	flag.StringVar(&rootArgs.database.password, "staple-db-password", "password123", "--gaia-bot-db-password password123")
	flag.StringVar(&rootArgs.hookSecret, "hook-secret", "", "--hook-secret asdf")
	flag.Parse()
}

func main() {

}
