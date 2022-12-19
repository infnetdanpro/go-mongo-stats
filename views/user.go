package views

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/infnetdanpro/go-mongo-stats/model"
	"github.com/infnetdanpro/go-mongo-stats/store"
)

func UserRegister(w http.ResponseWriter, r *http.Request, userRepo store.UserRepository) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user model.UserInput

	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(user.Password) < 3 {
		http.Error(w, "Password lenght must be greater than 3 symbols", http.StatusUnprocessableEntity)
		return
	}

	newUser, err := userRepo.New(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func UserLogin(w http.ResponseWriter, r *http.Request, cookieStore *sessions.CookieStore, userRepo store.UserRepository) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user model.UserInput
	err := decoder.Decode(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := cookieStore.Get(r, os.Getenv("COOKIE_NAME"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = userRepo.CheckPassword(user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	err = sessions.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDB, err := userRepo.GetByEmail(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get user by email
	userDB.Authenticated = true
	session.Values["user"] = userDB

	err = sessions.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userDB)
}

func UserCheck(w http.ResponseWriter, r *http.Request, cookieStore *sessions.CookieStore, userRepo store.UserRepository) {
	if r.Method != "GET" {
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

	userDB, err := userRepo.GetById(user.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(userDB)
}
