package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/infnetdanpro/go-mongo-stats/middleware"
	"github.com/infnetdanpro/go-mongo-stats/store"
	"github.com/infnetdanpro/go-mongo-stats/views"
)

type Server struct {
	AppRepository          store.AppRepository
	EventRepository        store.EventRepository
	StorageEventRepository store.EventStorageRepository
	UserRepository         store.UserRepository
	CookieStore            *sessions.CookieStore
}

func (s Server) Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", middleware.SetHTMLHeader(s.HomeHandler))
	r.HandleFunc("/api/health-check/", middleware.SetJSONHeader(s.HealthCheckHandler))

	// Get application info, registration by user
	r.HandleFunc("/api/app/check/", middleware.SetJSONHeader(s.GetAppByKeyHandler))
	r.HandleFunc("/api/app/register/", middleware.SetJSONHeader(s.RegisterAppHandler))

	// Work with stats: save, get
	r.HandleFunc("/api/stats/save/", middleware.SetJSONHeader(s.SaveStats))
	r.HandleFunc("/api/stats/", middleware.SetJSONHeader(s.GetStats))
	r.HandleFunc("/api/stats/distinct-values/", middleware.SetJSONHeader(s.GetDistinctValuesByField))

	// Work with users: CRUD
	r.HandleFunc("/api/user/register/", middleware.SetJSONHeader(s.UserRegister))
	r.HandleFunc("/api/user/login/", middleware.SetJSONHeader(s.UserLogin))
	r.HandleFunc("/api/user/check/", middleware.SetJSONHeader(s.UserCheck))

	log.Fatal(http.ListenAndServe(":8080", r))
}

func (s Server) HomeHandler(w http.ResponseWriter, r *http.Request) {
	views.Home(w)
}

func (s Server) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	views.HealthCheck(w)
}

func (s Server) RegisterAppHandler(w http.ResponseWriter, r *http.Request) {
	views.RegisterApp(w, r, s.AppRepository, s.UserRepository, s.CookieStore)
}

func (s Server) GetAppByKeyHandler(w http.ResponseWriter, r *http.Request) {
	views.GetAppByKey(w, r, s.AppRepository)
}

func (s Server) SaveStats(w http.ResponseWriter, r *http.Request) {
	views.SaveStats(w, r, s.AppRepository, s.EventRepository)
}

func (s Server) GetStats(w http.ResponseWriter, r *http.Request) {
	views.GetStats(w, r, s.StorageEventRepository, s.CookieStore, s.UserRepository)
}

func (s Server) GetDistinctValuesByField(w http.ResponseWriter, r *http.Request) {
	views.GetDistinctValuesByField(w, r, s.StorageEventRepository, s.CookieStore, s.UserRepository)
}

func (s Server) UserLogin(w http.ResponseWriter, r *http.Request) {
	views.UserLogin(w, r, s.CookieStore, s.UserRepository)
}

func (s Server) UserCheck(w http.ResponseWriter, r *http.Request) {
	views.UserCheck(w, r, s.CookieStore, s.UserRepository)
}

func (s Server) UserRegister(w http.ResponseWriter, r *http.Request) {
	views.UserRegister(w, r, s.UserRepository)
}
