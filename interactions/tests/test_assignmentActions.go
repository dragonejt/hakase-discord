package interactions

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/interactions"
	"github.com/dragonejt/hakase-discord/settings"
	"github.com/stretchr/testify/assert"
)





func TestUpdateAssignment(t *testing.T){

	bot, err := discordgo.New(fmt.Sprintf("Bot %s", settings.DISCORD_BOT_TOKEN))
	if err != nil {
		slog.Error(fmt.Sprintf("error creating discord session: %s", err.Error()))
		return
	}
	bot.StateEnabled = true

	tests:= []struct {
		name  					string
		userPermissions 		int
		customID 				string
		mockAssignmentID 		string
		mockAssignmentError 	error
		expectedContent 		string
		expectedError			error
	}{
		{
			name:            "User without admin permission",
			userPermissions: 0,
			customID:        "updateAssignment_123",
			expectedContent: "admin permissions needed!",
			expectedError: 	 nil,
		},
		{
			name:            "Valid assignment update",
			userPermissions: discordgo.PermissionAdministrator,
			customID:        "updateAssignment_123",
			mockAssignmentID: "123",
			mockAssignmentError: nil,
			expectedContent: "update assignment", 
			expectedError: 	 nil,
		},
		{
			name:            "Assignment not found",
			userPermissions: discordgo.PermissionAdministrator,
			customID:        "updateAssignment_999",
			mockAssignmentID: "999",
			mockAssignmentError: errors.New("assignment not found"),
			expectedContent: "assignment not found",
			expectedError: 	 nil,
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
                        Permissions: int64(tt.userPermissions),
                    },
                    
                    Data: &discordgo.MessageComponentInteractionData{
                        CustomID: tt.customID,
                    },
				},
			}
						
			interactions.UpdateAssignment(bot, interaction)
			
			assert.NoError(t, err, "expected no error but got: %v", err)
		})
	}
}