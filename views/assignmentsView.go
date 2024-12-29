package views

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
)

func AssignmentsListView(interactionCreate *discordgo.InteractionCreate, assignments []clients.Assignment) *discordgo.MessageEmbed {
	embed := discordgo.MessageEmbed{
		Title:       "assignments",
		Description: fmt.Sprintf("%d assignments in course", len(assignments)),
		Author:      &discordgo.MessageEmbedAuthor{Name: interactionCreate.Member.User.Username, IconURL: interactionCreate.Member.User.AvatarURL("")},
		Color:       0x00ff00,
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

func AssignmentsListActions(interactionCreate *discordgo.InteractionCreate, assignments []clients.Assignment) *discordgo.ActionsRow {
	addAssignmentButton := discordgo.Button{
		Emoji: &discordgo.ComponentEmoji{
			Name: "📝",
		},
		Label:    "add assignment",
		Style:    discordgo.PrimaryButton,
		CustomID: "addAssignment",
	}

	return &discordgo.ActionsRow{Components: []discordgo.MessageComponent{addAssignmentButton}}

}
