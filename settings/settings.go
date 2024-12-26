package settings

import "os"

var ENV string = os.Getenv("ENV")
var DEBUG bool = ENV != "production"
var DISCORD_APP_ID string = os.Getenv("DISCORD_APP_ID")
var DISCORD_BOT_TOKEN string = os.Getenv("DISCORD_BOT_TOKEN")
var BACKEND_URL string = os.Getenv("BACKEND_URL")
var BACKEND_API_KEY string = os.Getenv("BACKEND_API_KEY")
var NATS_URL string = os.Getenv("NATS_URL")
var STREAM_NAME string = os.Getenv("STREAM_NAME")
