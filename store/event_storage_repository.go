package store

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventStorageRepository struct {
	MG *mongo.Client
}

func EchoMongo(mongoDsn string) bool {
	_, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDsn))

	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	return true
}

func (e EventStorageRepository) GetEventsByFilters(bsonMapFilters bson.M) []map[string]interface{} {
	coll := e.MG.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))

	cursor, err := coll.Find(context.TODO(), bsonMapFilters)
	if err != nil {
		panic(err)
	}

	var events []map[string]interface{}
	if err = cursor.All(context.TODO(), &events); err != nil {
		panic(err)
	}
	return events
}

func (e EventStorageRepository) GetDistinctValuesByField(field string, bsonMapFilters bson.M) []interface{} {
	coll := e.MG.Database(os.Getenv("MONGO_DB")).Collection(os.Getenv("MONGO_COLLECTION"))

	results, err := coll.Distinct(context.TODO(), field, bsonMapFilters)
	if err != nil {
		panic(err)
	}

	return results
}

func (e EventStorageRepository) GetKeys(m map[string]interface{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
