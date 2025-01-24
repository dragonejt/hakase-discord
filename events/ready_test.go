package events

import (
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"github.com/dragonejt/hakase-discord/clients"
)

type MockSession struct {
	mock.Mock
}

func (m *MockSession) UpdateCustomStatus(status string) error {
	args := m.Called(status)
	return args.Error(0)
}

func (m *MockSession) UpdateCustomStatus(status string) error {
	args := m.Called(status)
	return args.Error(0)
}

func TestReady(t *testing.T) {
	mockSession := new(MockSession)
	mockReady := &discordgo.Ready{
		User: &discordgo.User{
			Username: "testuser",
		},
	}

	mockSession.On("UpdateCustomStatus", "assisting 0 classes").Return(nil)

	Ready(mockSession, mockReady)

	mockSession.AssertExpectations(t)
}

func TestReadyClientAPICall(t *testing.T) {
	mockSession := new(MockSession)
	mockReady := &discordgo.Ready{
		User: &discordgo.User{
			Username: "testuser",
		},
	}

	mockSession.On("UpdateCustomStatus", "assisting 0 classes").Return(nil)

	Ready(mockSession, mockReady)

	mockSession.AssertExpectations(t)

	// Assert that the correct clients API calls were called
	// This is a placeholder, replace with actual assertions based on your implementation
	// For example, if you have a mock client, you can assert that the expected methods were called
	// mockClient.AssertCalled(t, "ExpectedMethod", expectedArguments)
}
