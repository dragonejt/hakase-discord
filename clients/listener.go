package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dragonejt/hakase-discord/settings"
	"github.com/getsentry/sentry-go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/palantir/stacktrace"
)

func ListenToStream(bot *discordgo.Session, hakaseClient HakaseClient, stopListener chan bool) {
	slog.Info(fmt.Sprintf("opening NATS consumer connection to: %s", settings.NATS_URL))
	connection, err := nats.Connect(settings.NATS_URL)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error connecting to NATS: %s", settings.NATS_URL).Error())
		return
	}
	defer func() {
		slog.Info("draining NATS consumer connection")
		err := connection.Drain()
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "error draining NATS connection").Error())
			return
		}
	}()

	slog.Debug("opening jetstream consumer connection")
	js, err := jetstream.New(connection)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error opening jetstream publisher connection").Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Debug(fmt.Sprintf("creating stream with name: %s", settings.STREAM_NAME))
	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     settings.STREAM_NAME,
		Subjects: []string{fmt.Sprintf("%s.*", settings.STREAM_NAME)},
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error creating stream with name: %s", settings.STREAM_NAME).Error())
		return
	}

	consumer, err := js.CreateOrUpdateConsumer(ctx, settings.STREAM_NAME, jetstream.ConsumerConfig{
		Name:      settings.STREAM_NAME,
		Durable:   settings.STREAM_NAME,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error creating consumer for stream: %s", settings.STREAM_NAME).Error())
		return
	}

	subscription, err := consumer.Consume(func(message jetstream.Msg) {
		consumeMessage(bot, hakaseClient, message)
	})
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error subscribing to stream: %s", settings.STREAM_NAME).Error())
		return
	}
	defer subscription.Drain()

	<-stopListener

}

func consumeMessage(bot *discordgo.Session, hakaseClient HakaseClient, message jetstream.Msg) {
	transaction := sentry.StartTransaction(context.WithValue(context.Background(), DiscordSession{}, bot), "consumeMessage")
	defer transaction.Finish()
	slog.Info(fmt.Sprintf("received message: %s with subject: %s", string(message.Data()), message.Subject()))

	if message.Subject() == "notifications" {
		consumeNotification(transaction, hakaseClient, message)
	} else if message.Subject() == "assignments" {
		consumeAssignmentNotification(transaction, hakaseClient, message)
	} else {

	}

}

func consumeNotification(span *sentry.Span, hakaseClient HakaseClient, message jetstream.Msg) {
	span = span.StartChild("consumeNotification")
	defer span.Finish()

	slog.Info(fmt.Sprintf("received notification with message: %s", string(message.Data())))
}

func consumeAssignmentNotification(span *sentry.Span, hakaseClient HakaseClient, message jetstream.Msg) {
	span = span.StartChild("consumeAssignmentNotification")
	defer span.Finish()
	bot := span.GetTransaction().Context().Value(DiscordSession{}).(*discordgo.Session)

	assignmentNotification := AssignmentNotification{}
	err := json.Unmarshal(message.Data(), &assignmentNotification)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error unmarshalling assignment notification").Error())
		return
	}

	assignment, err := hakaseClient.ReadAssignment(span, fmt.Sprint(assignmentNotification.AssignmentID))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to get assignment with ID: %s", assignmentNotification.AssignmentID).Error())
		return
	}

	course, err := hakaseClient.ReadCourse(span, assignmentNotification.CourseID)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "failed to get course with courseID: %s", assignmentNotification.CourseID).Error())
		return
	}

	notificationsChannel := course.NotifyChannel
	if notificationsChannel == "" {
		guild, err := bot.Guild(course.CourseID)
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "unable to get guild system channel for notifications").Error())
			return
		}
		notificationsChannel = guild.SystemChannelID
	}

	notificationTime := assignment.Due.Add(-1 * assignmentNotification.Before)
	if time.Now().After(notificationTime) {
		_, err := bot.ChannelMessageSend(notificationsChannel, fmt.Sprintf("**[assignment notification]** assignment: %s is due in %s hours!", assignment.Name, assignmentNotification.Before/time.Hour))
		if err != nil {
			slog.Error(stacktrace.Propagate(err, "failed to send assignment notification for %s", assignment.ID).Error())
			// retry sending assignment notification in 15 minutes
			message.NakWithDelay(15 * time.Minute)
		}
	} else {
		err := message.NakWithDelay(time.Until(notificationTime))
		if err != nil {
			bot.ChannelMessageSend(notificationsChannel, fmt.Sprintf("**[assignment notification error]** failed to schedule assignment notifications for assignment: %s", assignment.Name))
			slog.Error(stacktrace.Propagate(err, "failed to schedule assignment notifications for %s", assignment.ID).Error())
		}
	}
}
