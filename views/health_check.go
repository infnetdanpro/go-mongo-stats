package views

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/infnetdanpro/go-mongo-stats/model"
	"github.com/infnetdanpro/go-mongo-stats/store"
)

func HealthCheck(w http.ResponseWriter) {
	var h model.HealthCheck

	h.Mongo = store.EchoMongo(os.Getenv("MONGODB_URI"))
	h.Rabbit = store.EchoRabbitMQ(os.Getenv("RABBITMQ_URL"))

	json.NewEncoder(w).Encode(h)
}
