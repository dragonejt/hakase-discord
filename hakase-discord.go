package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/dragonejt/hakase-discord/notifications"
	"github.com/getsentry/sentry-go"
)

func main() {
	debug := os.Getenv("ENV") != "production"
	if !debug {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:                "https://701b4c4b27e3aeb9ab991b282df7c705@o4507124907638784.ingest.us.sentry.io/4508476074360832",
			TracesSampleRate:   0.1,
			ProfilesSampleRate: 1,
			EnableTracing:      true,
			Environment:        os.Getenv("ENV"),
		})
		if err != nil {
			slog.Warn("error initiating sentry: " + err.Error())
		}

	} else {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	bot, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		slog.Error("error creating discord bot session" + err.Error())
		return
	}

	stopListener := make(chan bool, 1)
	go notifications.ListenToStream(stopListener)

	bot.AddHandler(events.Ready)

	err = bot.Open()
	if err != nil {
		slog.Error("error opening discord bot connection" + err.Error())
		return
	}
	defer bot.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	stopListener <- true
}
