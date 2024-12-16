package notifications

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartListening(rabbitMQURL string, queueName string) {
	connection, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		slog.Error("error starting rabbitMQ connection: " + err.Error())
		return
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		slog.Error("error initializing connection channel: " + err.Error())
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		slog.Error("error creating channel queue: " + err.Error())
	}

	slog.Info("listening to queue: " + queueName)
	notifications, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		slog.Error("error consuming notifications: " + err.Error())
	}

	for n := range notifications {
		slog.Info("received notification" + string(n.Body))
	}

}
