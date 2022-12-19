package store

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventRepository struct {
	Connection *amqp.Connection
}

func EchoRabbitMQ(rabbitDsn string) bool {
	conn, err := amqp.Dial(rabbitDsn)

	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	defer ch.Close()
	return true
}

// Save - save event to the RabbitMQ queue, then Consumer will save it to Mongo DB
func (e EventRepository) Save(data []byte, queueName string) (bool, error) {
	ch, err := e.Connection.Channel()
	if err != nil {
		if err != nil {
			log.Fatal(err.Error())
			return false, err
		}
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		if err != nil {
			log.Fatal(err.Error())
			return false, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		})

	if err != nil {
		log.Fatal(err.Error())
		return false, err
	}
	return true, nil
}
