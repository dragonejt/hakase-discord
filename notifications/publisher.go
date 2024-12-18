package notifications

import (
	"errors"
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
	slog.Debug("opening NATS connection to: " + natsURL)
	connection, err := nats.Connect(natsURL)
	if err != nil {
		panic(err.Error())
	}

	slog.Debug("opening jetstream connection")
	jsctx, err := connection.JetStream()
	if err != nil {
		panic(err.Error())
	}
	return jsctx
}

func PublishToStream(message string) error {
	streamName := os.Getenv("STREAM_NAME")
	jsctx := StreamConnectionPool.Get().(nats.JetStreamContext)
	defer StreamConnectionPool.Put(jsctx)

	slog.Info("publishing")
	_, err := jsctx.Publish(streamName+".notifications", []byte(message))
	if err != nil {
		return errors.New("error publishing to NATS subject: " + err.Error())
	}

	return nil
}
