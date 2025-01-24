package views

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func AssignmentsListView(member *discordgo.Member, assignments []clients.Assignment) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title:       "assignments",
		Description: fmt.Sprintf("%d assignments in course", len(assignments)),
		Author:      &discordgo.MessageEmbedAuthor{Name: member.User.Username, IconURL: member.User.AvatarURL("")},
	}

	for _, assignment := range assignments {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("id: %d", assignment.ID),
			Value:  assignment.Name,
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "due:",
			Value:  assignment.Due.Format(time.RFC1123),
			Inline: true,
		})
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "\u200B",
			Value:  "\u200B",
			Inline: false,
		})
	}

	return &embed

}

func AssignmentsListActions() *discordgo.ActionsRow {
	return &discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Emoji: &discordgo.ComponentEmoji{
					Name: "➕",
				},
				Label:    "add",
				Style:    discordgo.PrimaryButton,
				CustomID: "addAssignmentAction",
			},
		},
	}

}
