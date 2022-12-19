package views

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/infnetdanpro/go-mongo-stats/model"
	"github.com/infnetdanpro/go-mongo-stats/store"
)

func RegisterApp(w http.ResponseWriter, r *http.Request, appRepo store.AppRepository, userRepo store.UserRepository, cookieStore *sessions.CookieStore) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	session, err := cookieStore.Get(r, os.Getenv("COOKIE_NAME"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user := userRepo.GetUserFromSession(session)

	if auth := user.Authenticated; !auth {
		http.Error(w, "You must be authorizied", http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var app model.AppRegister

	err = decoder.Decode(&app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if app.Name == "" {
		http.Error(w, "'name' is required", http.StatusUnprocessableEntity)
		return
	}

	createdApp, err := appRepo.NewApp(app, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	json.NewEncoder(w).Encode(createdApp)
}

func GetAppByKey(w http.ResponseWriter, r *http.Request, appRepo store.AppRepository) bool {
	appKey := r.Header.Get("API-MaxPanel")

	if appKey == "" {
		http.Error(w, "API-MaxPanel is mandatory field", http.StatusUnprocessableEntity)
		return true
	}

	app, err := appRepo.GetAppByKey(appKey)

	if err != nil {
		http.Error(w, "App not found", http.StatusNotFound)
		return true
	}

	json.NewEncoder(w).Encode(app)
	return false
}
