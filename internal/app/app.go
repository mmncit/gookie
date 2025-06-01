package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	api "github.com/mmncit/gookie/internal/api/handlers"
)

type Application struct {
	Logger *log.Logger
	EventHandler *api.EventHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// stores
	

	// handlers
	eventHandler := api.NewEventHandler()

	app := &Application{
		Logger: logger,
		EventHandler: eventHandler,
	}
	return app, nil
}

func (app *Application)HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}