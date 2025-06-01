package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mmncit/gookie/internal/app"
	"github.com/mmncit/gookie/internal/routes"
)


func main() {

	var port int
	flag.IntVar(&port, "port", 8080, "Port to run the server on")
	flag.Parse()

	
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	app.Logger.Println("Application started successfully")


	
	routes := routes.SetupRoutes(app)
	
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: routes,
		IdleTimeout: time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("Starting server on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("Failed to start server: %v", err)
	} else {
		app.Logger.Println("Server is running on :8080")
	}
}

