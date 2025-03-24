package events_test

import (
	"log/slog"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GuildCreateTestSuite struct {
	suite.Suite
	bot          *discordgo.Session
	guildCreate  *discordgo.GuildCreate
	hakaseClient *MockHakaseClient
}

type MockHakaseClient struct {
	clients.HakaseClient
	mock.Mock
}

func (hakaseClient *MockHakaseClient) CreateCourse(span *sentry.Span, course clients.Course) error {
	hakaseClient.Called(span, course)
	return nil
}

func (testSuite *GuildCreateTestSuite) SetupTest() {
	testSuite.bot = &discordgo.Session{}
	testSuite.bot.State = &discordgo.State{}
	testSuite.bot.State.Guilds = []*discordgo.Guild{}
	testSuite.guildCreate = &discordgo.GuildCreate{
		Guild: &discordgo.Guild{
			ID:   "1234567890",
			Name: "Test Guild",
		},
	}
	testSuite.hakaseClient = new(MockHakaseClient)
}

func (testSuite *GuildCreateTestSuite) TestGuildCreateSuccess() {
	testSuite.hakaseClient.On("CreateCourse", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info("CreateCourse called")
	})
	events.GuildCreate(testSuite.bot, testSuite.guildCreate, testSuite.hakaseClient)
	testSuite.hakaseClient.AssertExpectations(testSuite.T())
}

func TestGuildEvents(t *testing.T) {
	suite.Run(t, new(GuildCreateTestSuite))
}
