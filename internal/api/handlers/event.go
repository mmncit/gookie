package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
   
}

func NewEventHandler() *EventHandler {
    return &EventHandler{}
}

// GET /events/{id}
func (eh *EventHandler) HandleGetEventByID(w http.ResponseWriter, r *http.Request) {
   paramsEventId := chi.URLParam(r, "id")
   if paramsEventId == "" {
	  http.Error(w, "Event ID is required", http.StatusBadRequest)
	  return
   }
   eventID , err := strconv.ParseInt(paramsEventId, 10, 64)
   if err != nil {
	  http.Error(w, "Invalid Event ID", http.StatusBadRequest)
	  return
   }
   fmt.Fprintf(w, "Event ID: %d", eventID)
}

// POST /events
func (eh *EventHandler) HandleCreateEvent(w http.ResponseWriter, r *http.Request) {
   // Here you would typically parse the request body to create a new event
   // For simplicity, we will just return a success message
   w.WriteHeader(http.StatusCreated)
   fmt.Fprintln(w, "Event created successfully")
}

