package notifications

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

func StartListening(rabbitMQURL string, queueName string) error {
	connection, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return err
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	notifications, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for n := range notifications {
		slog.Info("received notification" + string(n.Body))
	}

	return nil

}
