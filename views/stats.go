package views

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/infnetdanpro/go-mongo-stats/model"
	"github.com/infnetdanpro/go-mongo-stats/store"
)

func SaveStats(w http.ResponseWriter, r *http.Request, appRepo store.AppRepository, eventRepo store.EventRepository) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	appKey := r.Header.Get("API-MaxPanel")
	if appKey == "" {
		http.Error(w, "API-MaxPanel is mandatory header", http.StatusUnprocessableEntity)
		return
	}

	app, err := appRepo.GetAppByKey(appKey)

	if err != nil {
		http.Error(w, "App not found", http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	var event map[string]interface{}

	json.Unmarshal(body, &event)
	event["__timestamp"] = time.Unix(time.Now().Unix(), 0)
	event["__app_id"] = app.Id
	event["__user_agent"] = r.Header.Get("user-agent")
	event["__ip"] = r.RemoteAddr

	if err != nil {
		http.Error(w, "Wrong json, sorry", http.StatusBadRequest)
		return
	}

	eventData, err := json.Marshal(event)
	if err != nil {
		http.Error(w, "Wrong event json parser, sorry", http.StatusBadRequest)
		return
	}
	_, err = eventRepo.Save(eventData, os.Getenv("QUEUE_NAME"))
	if err != nil {
		http.Error(w, "Problem with rabbitmq saving", http.StatusInternalServerError)
	}
}

func GetStats(w http.ResponseWriter, r *http.Request, eventStorageRepo store.EventStorageRepository, cookieStore *sessions.CookieStore, userRepo store.UserRepository) {
	// POST BECAUSE WE NEED TO CREATE A BIG FILTER
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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var filters model.Filters
	json.Unmarshal(body, &filters)

	events := eventStorageRepo.GetEventsByFilters(filters.Filters)

	for _, event := range events {
		keys := eventStorageRepo.GetKeys(event)
		for _, key := range keys {
			value := event[key]
			if event[key] != nil && strings.HasPrefix(key, "__") {
				replacedKey := strings.ReplaceAll(key, "__", "")
				event[replacedKey] = value
				delete(event, key)
			}
		}
	}

	json.NewEncoder(w).Encode(events)
}

func GetDistinctValuesByField(w http.ResponseWriter, r *http.Request, eventStorageRepo store.EventStorageRepository, cookieStore *sessions.CookieStore, userRepo store.UserRepository) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var filters model.DistinctFilters
	json.Unmarshal(body, &filters)

	events := eventStorageRepo.GetDistinctValuesByField(filters.Field, filters.Filters)
	json.NewEncoder(w).Encode(events)

}
