package routes

import (
	"net/http"
	"screen-therapy-backend/handlers"

	"github.com/gorilla/mux"
)

// Test handler function
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🏡 Welcome to the Screen Therapy API!"))
}

// RegisterRoutes sets up API endpoints
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// 🔥 User Authentication Routes
	r.HandleFunc("/checkUser/{userId}", handlers.CheckUser).Methods("GET") // 🔍 Check if a user exists
	r.HandleFunc("/registerUser", handlers.RegisterUser).Methods("POST")   // 🆕 Register a new user
}
