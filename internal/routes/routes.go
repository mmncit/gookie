package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmncit/gookie/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	routes := chi.NewRouter()
	routes.Get("/health", app.HealthCheckHandler)
	return routes
	
}