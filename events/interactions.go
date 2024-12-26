package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/commands"
)

func InteractionCreate(bot *discordgo.Session, interactionCreate *discordgo.InteractionCreate) {
	switch interactionCreate.ApplicationCommandData().Name {
	case "assignments":
		commands.SlashAssignments(bot, interactionCreate)
	}
}
