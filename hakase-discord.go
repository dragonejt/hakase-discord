package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/commands"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/dragonejt/hakase-discord/notifications"
	"github.com/dragonejt/hakase-discord/settings"
	"github.com/getsentry/sentry-go"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if !settings.DEBUG {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:                "https://701b4c4b27e3aeb9ab991b282df7c705@o4507124907638784.ingest.us.sentry.io/4508476074360832",
			TracesSampleRate:   0.1,
			ProfilesSampleRate: 1,
			EnableTracing:      true,
			Environment:        settings.ENV,
		})
		if err != nil {
			slog.Warn(fmt.Sprintf("error initiating sentry: %s", err))
		}

	} else {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", settings.DISCORD_BOT_TOKEN))
	if err != nil {
		slog.Error(fmt.Sprintf("error creating discord session: %s", err.Error()))
		return
	}
	bot.StateEnabled = true

	stopListener := make(chan bool, 1)
	go notifications.ListenToStream(stopListener)

	slog.Info("registering event handlers")
	bot.AddHandler(events.Ready)
	bot.AddHandler(events.GuildCreate)
	bot.AddHandler(events.GuildDelete)
	bot.AddHandler(events.InteractionCreate)

	slog.Info("registering commands")
	commands := []*discordgo.ApplicationCommand{commands.AssignmentsCommand}
	for _, cmd := range commands {
		_, err = bot.ApplicationCommandCreate(settings.DISCORD_APP_ID, "", cmd)
		if err != nil {
			slog.Error(fmt.Sprintf("error registering command: %s, %s", cmd.Name, err.Error()))
		} else {
			slog.Info(fmt.Sprintf("successfully registered command: %s", cmd.Name))
		}
	}

	err = bot.Open()
	if err != nil {
		slog.Error(fmt.Sprintf("error opening connection to discord: %s", err.Error()))
		return
	}
	defer bot.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	stopListener <- true
}
