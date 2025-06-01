package main

import (
	"net/http"
	"time"

	"github.com/mmncit/gookie/internal/app"
)


func main() {

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	app.Logger.Println("Application started successfully")
	server := &http.Server{
		Addr:    ":8080",
		IdleTimeout: time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("Failed to start server: %v", err)
	} else {
		app.Logger.Println("Server is running on :8080")
	}
}

