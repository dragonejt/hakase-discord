package events

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Ready(session *discordgo.Session, ready *discordgo.Ready, channel *amqp.Channel) {
	slog.Info("logged in as " + ready.User.String())
}
