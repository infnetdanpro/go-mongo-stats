package worker

import (
	"context"
	"encoding/json"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Recieve message from RabbitMQ and save to the Mongo DB
func QueueConsumer(conn *amqp.Connection, mgClient *mongo.Client, queueName string) {
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	collection := mgClient.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			// log.Printf("Received a message: %s", d.Body)
			var jsonData interface{}
			json.Unmarshal(d.Body, &jsonData)

			_, err := collection.InsertOne(context.TODO(), jsonData)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Println("Inserted a single document: ", insertResult.InsertedID)
		}
	}()

	log.Printf(" [*] Waiting for messages.")
	<-forever
}
