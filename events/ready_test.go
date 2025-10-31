package events_test

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/clients"
	"github.com/dragonejt/hakase-discord/events"
	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ReadyEventsTestSuite struct {
	suite.Suite
	bot           *discordgo.Session
	ready         *discordgo.Ready
	notifications *MockNotificationsClient
	hakaseClient  clients.HakaseClient
}

func TestReadyEvents(t *testing.T) {
	suite.Run(t, new(ReadyEventsTestSuite))
}

type MockNotificationsClient struct {
	clients.NotificationsClient
	mock.Mock
}

func (notifications *MockNotificationsClient) PublishNotification(span *sentry.Span, notification string) {
	notifications.Called(span, notification)
}

func (testSuite *ReadyEventsTestSuite) SetupTest() {
	testSuite.bot = &discordgo.Session{}
	testSuite.bot.State = &discordgo.State{}
	testSuite.ready = &discordgo.Ready{
		User: &discordgo.User{
			ID:       "1234567890",
			Username: "Test User",
		},
	}
	testSuite.notifications = new(MockNotificationsClient)
	testSuite.hakaseClient = clients.HakaseClient{
		Backend:       nil,
		Notifications: testSuite.notifications,
	}
}

func (testSuite *ReadyEventsTestSuite) TestReady() {
	testSuite.notifications.On("PublishNotification", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		slog.Info(fmt.Sprintf("PublishNotification called with: %s", args.String(1)))
	})
	events.Ready(testSuite.bot, testSuite.ready, testSuite.hakaseClient)
	testSuite.notifications.AssertExpectations(testSuite.T())
}
