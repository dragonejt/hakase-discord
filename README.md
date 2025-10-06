# hakase-discord
[![codecov](https://codecov.io/gh/dragonejt/hakase-discord/graph/badge.svg?token=7MEF3IHI00)](https://codecov.io/gh/dragonejt/hakase-discord)

hakase is a collection of helpful utilities for class chatrooms, including an assignment due date reminder, study session scheduler, and more. It is currently under development. This repository holds the Discord Bot, built with Go. The backend API is in a beta state, with a Slack App planned in the future.

Backend API: https://github.com/dragonejt/hakase

## Table of Contents
- Local Development
  - Building and Running
  - Testing
  - Linting and Formatting
- Deployment
  - Continuous Delivery

## Local Development
### Building and Running
Local development with hakase-discord is relatively simple. The only command you have to run is:
```sh
go run hakase-discord.go
```
You do have to have some environment variables in place. hakase does not directly read from a .env file, but you can configure environment variables or reference a .env file through IDE launch options. Otherwise, you can set environment variables locally.
```sh
ENV="development"
DISCORD_BOT_TOKEN="from Discord Dev Portal"
BACKEND_URL="https://hakase.dragonejt.dev" # if self-hosting, change URL to self-hosted backend
BACKEND_API_KEY="from Backend API"
NATS_URL="nats://"
STREAM_NAME="hakase_discord_local" # different from production stream name
```

### Testing
Testing has not been implemented for hakase-discord yet.

### Linting and Formatting
Go and the [VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go) automatically performs linting and formatting on save. 
The `integrate.yml` GitHub Actions workflow will check for linting errors and formatting mistakes with [golangci-lint](https://github.com/golangci/golangci-lint-action).

## Deployment
For deployment, hakase is built into a Docker image with [nixpacks](https://nixpacks.com/docs/providers/go), and then deployed into a container via [Dokku](https://dokku.com/).

On the deployed docker container, the following environment variables should be set:
```sh
ENV="production"
DISCORD_BOT_TOKEN="from Discord Dev Portal"
BACKEND_URL="https://hakase.dragonejt.dev" # if self-hosting, change URL to self-hosted backend
BACKEND_API_KEY="from Backend API"
NATS_URL="nats://"
STREAM_NAME="hakase_discord" # different from development stream name
SENTRY_DSN="from Sentry"
```
Dokku does support dockerized message queues, and hakase uses a dockerized NATS instance in production.

### Continuous Delivery
hakase has a continuous delivery GitHub Actions workflow, `deliver.yml`. The steps taken are summarized:

1. Build a Docker image with the [nixpacks GitHub Action](https://github.com/iloveitaly/github-action-nixpacks)
2. The nixpacks GitHub Action uploads the built image to [GitHub Container Registry](https://github.com/dragonejt/hakase-discord/pkgs/container/hakase-discord)
3. The built docker image is deployed as a docker container via the [Dokku GitHub Action](https://github.com/dokku/github-action)
4. A new Sentry release is created for monitoring with the [Sentry Release GitHub Action](https://github.com/getsentry/action-release).