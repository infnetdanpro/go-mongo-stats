package main

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"

	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/infnetdanpro/go-mongo-stats/store"
	"github.com/joho/godotenv"
)

//go:embed migrations
var migrations embed.FS

const schemaVersion = 1

func ensureSchema(db *sql.DB) error {
	sourceInstance, err := httpfs.New(http.FS(migrations), os.Getenv("MIGRATION_DIR"))
	if err != nil {
		return fmt.Errorf("invalid source instance, %w", err)
	}
	targetInstance, err := sqlite3.WithInstance(db, new(sqlite3.Config))
	if err != nil {
		return fmt.Errorf("invalid target sqlite instance, %w", err)
	}
	m, err := migrate.NewWithInstance("httpfs", sourceInstance, "sqlite3", targetInstance)
	if err != nil {
		return fmt.Errorf("failed to initialize migrate instance, %w", err)
	}
	err = m.Migrate(schemaVersion)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Migration applied")
	return sourceInstance.Close()
}
func Goro() {
	log.Println(11111111)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))

	if err != nil {
		log.Fatal(err.Error())
		panic("Problem with DB!")
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err.Error())
		panic(err.Error())
	}

	// Apply migrations
	migrationError := ensureSchema(db)

	if migrationError != nil {
		log.Fatal(migrationError.Error())
		panic(migrationError.Error())
	}

	appRepo := store.AppRepository{DB: db}
	// eventRepo := store.EventRepository{}
	// server := Server{AppRepository: appRepo, EventRepository: eventRepo}
	server := Server{AppRepository: appRepo}

	log.Println("Ready to start")
	go Goro()
	go server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
