package store

import (
	"database/sql"
	"log"

	"github.com/gorilla/sessions"
	"github.com/infnetdanpro/go-mongo-stats/model"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	DB            *sql.DB
	CookieSession *sessions.Session
}

func (u UserRepository) New(email string, password string) (*model.User, error) {
	user := &model.User{}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Problem with user password hash!")
		return &model.User{}, err
	}
	err = u.DB.QueryRow("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, is_active", email, string(hash)).Scan(&user.ID, &user.Email, &user.IsActive)
	if err != nil {
		return &model.User{}, err
	}
	return user, nil
}

func (u UserRepository) CheckPassword(email string, password string) error {
	userDB := &model.CheckUser{}

	err := u.DB.QueryRow("SELECT password FROM users WHERE email = $1  AND is_active = 1", email).Scan(&userDB.Password)
	if err != nil {
		return err
	}

	errr := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(password))

	if errr != nil {
		log.Fatal("Problem with user password hash!")
		return err
	}

	return nil
}

func (u UserRepository) GetByEmail(email string) (*model.User, error) {
	user := &model.User{}

	err := u.DB.QueryRow("SELECT id, email, is_active FROM users WHERE email = $1 AND is_active = 1", email).Scan(&user.ID, &user.Email, &user.IsActive)
	if err != nil {
		return &model.User{}, err
	}
	return user, nil
}

func (u UserRepository) GetUserFromSession(s *sessions.Session) *model.User {
	val := s.Values["user"]
	var user = model.User{}

	user, ok := val.(model.User)

	if !ok {
		return &model.User{Authenticated: false}
	}
	user.Authenticated = true
	return &user
}

func (u UserRepository) GetById(userId int) (*model.User, error) {
	user := &model.User{}
	err := u.DB.QueryRow("SELECT id, email, is_active, TRUE as 'authenticated' FROM users WHERE id = $1 and is_active = 1", userId).Scan(&user.ID, &user.Email, &user.IsActive, &user.Authenticated)
	if err != nil {
		return &model.User{}, err
	}
	return user, nil
}
