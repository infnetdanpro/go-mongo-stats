package model

type App struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Key      string `json:"key"`
	IsActive bool   `json:"is_active"`
}

type AppRegister struct {
	Name string `json:"name"`
}
