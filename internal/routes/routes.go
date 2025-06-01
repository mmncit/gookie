package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmncit/gookie/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	routes := chi.NewRouter()

	// health check route
	routes.Get("/health", app.HealthCheckHandler)

	// Event routes
	routes.Get("/events/{id}", app.EventHandler.HandleGetEventByID)
	routes.Post("/events", app.EventHandler.HandleCreateEvent)

	return routes
	
}