package store

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func GetDB(dbDriver string, dbName string) *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")

	if err != nil {
		log.Fatal(err.Error())
	}

	if err = db.Ping(); err != nil {
		panic(err.Error())
	}
	return db
}
