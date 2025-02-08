package interactions

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/interactions"
)



func testUpdateAssignment(t *testing.T){

	bot := &discordgo.Session{}

	tests:= []struct {
		name  string
		userPermissions int
		customID string
		expectedContent string
	}{
		{
			name:            "User without admin permission",
			userPermissions: 0,
			customID:        "updateAssignment_123",
			expectedContent: "admin permissions needed!",
		},
		{
			name:            "Valid assignment update",
			userPermissions: discordgo.PermissionAdministrator,
			customID:        "updateAssignment_123",
			expectedContent: "update assignment", 
		},
		{
			name:            "Assignment not found",
			userPermissions: discordgo.PermissionAdministrator,
			customID:        "updateAssignment_999",
			expectedContent: "assignment not found",
		},
	}

	for _, tt := range tests{
		t.Run(tt.name, func(t *testing.T){
			interaction := &discordgo.InteractionCreate{
				Interaction: &discordgo.Interaction{
					Member: &discordgo.Member{
                        User: &discordgo.User{
							Username: "TestUser", ID: "123456",
						},
                        Permissions: tt.userPermissions,
                    },
                    
                    Data: &discordgo.MessageComponentInteractionData{
                        CustomID: tt.customID,
                    },
				},

			}

			response := ""
			bot.InteractionRespond = func(i *discordgo.Interaction, ir *discordgo.InteractionResponse) error {
				if ir.Data!= nil{
					response = ir.Data.Content
				}
				return nil
			}

			interactions.UpdateAssignment(bot, interaction)
			assert.true(t, strings.Contains(response, tt.expectedContent), "expected response content to contain: %s, but got: %s", tt.expectedContent, response)
		})
	}
}