package store

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/infnetdanpro/go-mongo-stats/model"
)

type AppRepository struct {
	DB *sql.DB
}

func (a AppRepository) GetAppByKey(appKey string) (*model.App, error) {
	app := &model.App{}
	err := a.DB.QueryRow("SELECT id, name, key, is_active FROM apps WHERE key = $1", appKey).Scan(&app.Id, &app.Name, &app.Key, &app.IsActive)
	if err != nil {
		return &model.App{}, err
	}
	return app, nil
}

func (a AppRepository) NewApp(newApp model.AppRegister) (*model.App, error) {
	app := &model.App{}
	id := uuid.New()
	err := a.DB.QueryRow("INSERT INTO apps (name, key) VALUES ($1, $2) RETURNING id, name, key, is_active", newApp.Name, id.String()).Scan(&app.Id, &app.Name, &app.Key, &app.IsActive)

	if err != nil {
		return &model.App{}, err
	}
	return app, nil

}
