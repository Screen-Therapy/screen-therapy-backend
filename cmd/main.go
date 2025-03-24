package main

import (
	"fmt"
	"log"
	"net/http"

	"screen-therapy-backend/config"
	"screen-therapy-backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Firebase & Firestore
	config.InitFirebase()
	defer config.Client.Close()

	// Create a new router
	r := mux.NewRouter()

	// Register API routes
	routes.RegisterRoutes(r)

	// Start the HTTP server on localhost:8080
	fmt.Println("🚀 Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
