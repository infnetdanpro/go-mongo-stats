package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/infnetdanpro/go-mongo-stats/middleware"
	"github.com/infnetdanpro/go-mongo-stats/model"
	"github.com/infnetdanpro/go-mongo-stats/store"
)

type Server struct {
	AppRepository store.AppRepository
	// EventRepository store.EventRepository
}

func (s Server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.SetHTMLHeader(s.HomeHandler))
	r.HandleFunc("/api/health-check/", middleware.SetJSONHeader(s.HealthCheckHandler))
	r.HandleFunc("/api/check/", middleware.SetJSONHeader(s.GetAppByKeyHandler))
	r.HandleFunc("/api/register/", middleware.SetJSONHeader(s.RegisterAppHandler))
	r.HandleFunc("/api/stats/", middleware.SetJSONHeader(s.SaveStats))

	log.Fatal(http.ListenAndServe(":8088", r))
}

func (s Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(os.Getenv("TEMPLATED_DIR") + "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, "test")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	var h model.HealthCheck

	h.Mongo = store.EchoMongo(os.Getenv("MONGODB_URI"))
	h.Rabbit = store.EchoRabbitMQ(os.Getenv("RABBITMQ_URL"))

	json.NewEncoder(w).Encode(h)
}

func (s Server) RegisterAppHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var app model.AppRegister

	err := decoder.Decode(&app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if app.Name == "" {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	createdApp, err := s.AppRepository.NewApp(app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(createdApp)
}

func (s Server) GetAppByKeyHandler(w http.ResponseWriter, r *http.Request) {
	appKey := r.Header.Get("API-MaxPanel")

	if appKey == "" {
		http.Error(w, "API-MaxPanel is mandatory field", http.StatusUnprocessableEntity)
		return
	}

	app, err := s.AppRepository.GetAppByKey(appKey)

	if err != nil {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(app)
}

func (s Server) SaveStats(w http.ResponseWriter, r *http.Request) {
	// get connection to the rabbit and put data into queue
}
