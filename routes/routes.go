package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Test handler function
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🏡 Welcome to the Screen Therapy API!"))
}

// RegisterRoutes sets up API routes
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")
}
