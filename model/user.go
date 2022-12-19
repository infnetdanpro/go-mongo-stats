package model

type User struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	IsActive      bool   `json:"is_active"`
	Authenticated bool   `json:"authenticated"`
}

type CheckUser struct {
	Password string `json:"password"`
}

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
