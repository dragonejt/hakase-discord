// Package settings provides environment variable configuration for the bot.
package settings

import "os"

var ENV string = os.Getenv("ENV")
var DEBUG bool = ENV != "production"
var DISCORD_BOT_TOKEN string = os.Getenv("DISCORD_BOT_TOKEN")
var BACKEND_URL string = os.Getenv("BACKEND_URL")
var BACKEND_API_KEY string = os.Getenv("BACKEND_API_KEY")
var NATS_URL string = os.Getenv("NATS_URL")
var STREAM_NAME string = os.Getenv("STREAM_NAME")
var SENTRY_DSN string = os.Getenv("SENTRY_DSN")
