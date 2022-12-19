package main

import (
	"context"
	"database/sql"
	"embed"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/golang-migrate/migrate/v4/source/httpfs"
	"github.com/infnetdanpro/go-mongo-stats/model"
	"github.com/infnetdanpro/go-mongo-stats/store"
	"github.com/infnetdanpro/go-mongo-stats/worker"
	"github.com/joho/godotenv"
)

//go:embed migrations
var migrations embed.FS

// Go Migrations latest up version
const schemaVersion = 3

// Cookie Based authorization
var cookieStore *sessions.CookieStore

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

func initSessionStore() {
	// Prepare cookie based authorization for web
	cookieStore = sessions.NewCookieStore(
		securecookie.GenerateRandomKey(64),
		securecookie.GenerateRandomKey(32),
	)

	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 15,
		HttpOnly: true,
	}
	gob.Register(model.User{})
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
	defer db.Close()

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

	initSessionStore()

	qConnection, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		panic(err.Error())
	}
	defer qConnection.Close()

	mgClient, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))

	appRepo := store.AppRepository{DB: db}
	eventRepo := store.EventRepository{Connection: qConnection}
	eventStorageRepo := store.EventStorageRepository{MG: mgClient}
	userRepo := store.UserRepository{DB: db}
	server := Server{
		AppRepository:          appRepo,
		EventRepository:        eventRepo,
		StorageEventRepository: eventStorageRepo,
		UserRepository:         userRepo,
		CookieStore:            cookieStore,
	}

	log.Println("Ready to start")
	go worker.QueueConsumer(qConnection, mgClient, os.Getenv("QUEUE_NAME"))
	go server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
