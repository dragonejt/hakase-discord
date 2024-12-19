package notifications

import (
	"encoding/json"
	"log/slog"
	"os"
	"sync"

	"github.com/nats-io/nats.go"
)

var StreamConnectionPool sync.Pool = sync.Pool{
	New: createStreamConnection,
}

func createStreamConnection() any {
	natsURL := os.Getenv("NATS_URL")
	slog.Debug("opening NATS publisher connection to: " + natsURL)
	connection, err := nats.Connect(natsURL)
	if err != nil {
		panic(err.Error())
	}

	slog.Debug("opening jetstream publisher connection")
	jsctx, err := connection.JetStream()
	if err != nil {
		panic(err.Error())
	}
	return jsctx
}

func publishMessage(subject string, message []byte) error {
	streamName := os.Getenv("STREAM_NAME")
	jsctx := StreamConnectionPool.Get().(nats.JetStreamContext)
	defer StreamConnectionPool.Put(jsctx)

	slog.Info("publishing to subject: " + streamName + "." + subject)
	_, err := jsctx.Publish(streamName+"."+subject, message)
	if err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func PublishNotification(notification string) {
	err := publishMessage("notifications", []byte(notification))
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func PublishAssignmentNotification(notification AssignmentNotification) {
	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
	err = publishMessage("assignments", message)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func PublishStudySessionNotification(notification StudySessionNotification) {
	message, err := json.Marshal(notification)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	err = publishMessage("study_sessions", message)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
