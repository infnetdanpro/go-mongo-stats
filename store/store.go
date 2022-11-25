package store

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetDB(dbDriver string, dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		panic(err.Error())
	}
	return db
}

func EchoMongo(mongoDsn string) bool {
	_, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDsn))

	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	return true
}

func EchoRabbitMQ(rabbitDsn string) bool {
	conn, err := amqp.Dial(rabbitDsn)
	defer conn.Close()

	if err != nil {
		log.Fatal(err.Error())
		return false
	}

	ch, err := conn.Channel()
	defer ch.Close()

	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	return true
}
