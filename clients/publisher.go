package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dragonejt/hakase-discord/settings"
	"github.com/getsentry/sentry-go"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/palantir/stacktrace"
)

var PublisherPool sync.Pool = sync.Pool{
	New: createStreamConnection,
}

func createStreamConnection() any {
	slog.Info(fmt.Sprintf("opening NATS publisher connection to: %s", settings.NATS_URL))
	connection, err := nats.Connect(settings.NATS_URL)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error connecting to NATS: %s", settings.NATS_URL).Error())
		return nil
	}

	slog.Debug("opening jetstream publisher connection")
	js, err := jetstream.New(connection)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error opening jetstream publisher connection").Error())
		return nil
	}
	return js
}

func publishMessage(span *sentry.Span, subject string, message []byte) error {
	js := PublisherPool.Get().(jetstream.JetStream)
	defer PublisherPool.Put(js)

	ctx, cancel := context.WithTimeout(span.Context(), 10*time.Second)
	defer cancel()

	slog.Debug(fmt.Sprintf("publishing message to subject: %s.%s", settings.STREAM_NAME, subject))
	_, err := js.Publish(ctx, fmt.Sprintf("%s.%s", settings.STREAM_NAME, subject), message)
	if err != nil {
		return stacktrace.Propagate(err, "error publishing message to subject: %s.%s", settings.STREAM_NAME, subject)
	}

	return nil
}

func PublishNotification(span *sentry.Span, notification string) {
	span = span.StartChild("publishNotification")
	defer span.Finish()

	err := publishMessage(span, "notifications", []byte(notification))
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error publishing notification").Error())
		return
	}
}

func PublishAssignmentNotification(span *sentry.Span, notification AssignmentNotification) {
	span = span.StartChild("publishAssignmentNotification")
	defer span.Finish()

	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error marshalling assignment notification").Error())
		return
	}
	err = publishMessage(span, "assignments", message)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error publishing assignment notification").Error())
		return
	}
}

func PublishStudySessionNotification(span *sentry.Span, notification StudySessionNotification) {
	span = span.StartChild("publishStudySessionNotification")
	defer span.Finish()

	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error marshalling study session notification").Error())
		return
	}

	err = publishMessage(span, "study_sessions", message)
	if err != nil {
		slog.Error(stacktrace.Propagate(err, "error publishing study session notification").Error())
		return
	}
}
