package settings

import "os"

var ENV string = os.Getenv("ENV")
var DEBUG bool = ENV != "production"
var DISCORD_BOT_TOKEN string = os.Getenv("DISCORD_BOT_TOKEN")
var NATS_URL string = os.Getenv("NATS_URL")
var STREAM_NAME string = os.Getenv("STREAM_NAME")
