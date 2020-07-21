# gaia-bot

This is the bot taking care of new PRs and testing commands. This bot can be invoked by calling
`/test` on an open Gaia PR.

# Executors

As of this writing (2020.07.21) the bot uses a remote hetzner machine in order to build and push
docker images because it's too expensive to use machines in circleci.
