package views

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func ConfigView(course clients.Course) *discordgo.MessageEmbed {
	notifyChannel, notifyRole := course.NotifyChannel, course.NotifyRole
	if notifyChannel != "" {
		notifyChannel = fmt.Sprintf("<#%s>", notifyChannel)
	}
	if notifyRole != "" {
		notifyRole = fmt.Sprintf("<@&%s>", notifyRole)
	}

	return &discordgo.MessageEmbed{
		Title: "course config",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "notifications channel",
				Value: notifyChannel,
			},
			{
				Name:  "notifications role",
				Value: notifyRole,
			},
		},
	}
}

func ConfigActions() *discordgo.ActionsRow {
	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "📝",
				},
				Label:    "edit config",
				Style:    discordgo.PrimaryButton,
				CustomID: "configureCourse",
			},
		},
	}
}
