package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/dragonejt/hakase-discord/interactions"
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
			slog.Warn(fmt.Sprintf("Error Initiating Sentry: %s", err))
		}

	} else {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", settings.DISCORD_BOT_TOKEN))
	if err != nil {
		slog.Error(fmt.Sprintf("Error Creating Discord Session: %s", err.Error()))
		return
	}
	bot.StateEnabled = true

	stopListener := make(chan bool, 1)
	go notifications.ListenToStream(stopListener)

	err = bot.Open()
	if err != nil {
		slog.Error(fmt.Sprintf("Error Opening Connection to Discord: %s", err.Error()))
		return
	}
	defer bot.Close()

	slog.Info("Registering Event Handlers")
	bot.AddHandler(events.Ready)
	bot.AddHandler(events.GuildCreate)
	bot.AddHandler(events.GuildDelete)
	bot.AddHandler(events.InteractionCreate)

	slog.Info("Registering Interactions")
	interactions := []*discordgo.ApplicationCommand{&interactions.AssignmentsCommand, &interactions.HakaseCommand}
	for _, cmd := range interactions {
		_, err = bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd)
		if err != nil {
			slog.Error(fmt.Sprintf("Error Registering Command: %s, %s", cmd.Name, err.Error()))
		} else {
			slog.Info(fmt.Sprintf("Successfully Registered Command: %s", cmd.Name))
		}
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	stopListener <- true
}
