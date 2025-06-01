package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	api "github.com/mmncit/gookie/internal/api/handlers"
	"github.com/mmncit/gookie/store"
)

type Application struct {
	Logger *log.Logger
	EventHandler *api.EventHandler
	DB *sql.DB
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}


	// handlers
	eventHandler := api.NewEventHandler()

	app := &Application{
		Logger: logger,
		EventHandler: eventHandler,
		DB: pgDB,
	}
	return app, nil
}

func (app *Application)HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}