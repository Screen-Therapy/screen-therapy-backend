package routes

import (
	"net/http"
	"screen-therapy-backend/handlers"
	"screen-therapy-backend/handlers/auth"

	"github.com/gorilla/mux"
)

// Test handler function
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("🏡 Welcome to the Screen Therapy API!"))
}

// RegisterRoutes sets up API endpoints
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", HomeHandler).Methods("GET")

	// 🍏 Apple Authentication Routes
	r.HandleFunc("/apple/checkUser/{userId}", auth.CheckAppleUser).Methods("GET")
	r.HandleFunc("/apple/register", auth.RegisterAppleUser).Methods("POST")
	r.HandleFunc("/apple/checkUsername/{userId}", auth.CheckAppleUsername).Methods("GET")
	r.HandleFunc("/apple/setUsername", auth.SetAppleUsername).Methods("POST")

	// 📧 Email Auth Routes
	r.HandleFunc("/email/register", auth.RegisterEmailUser).Methods("POST")
	r.HandleFunc("/email/login", auth.LoginEmailUser).Methods("POST")

	// user routes
	r.HandleFunc("/user/info/{userId}", handlers.GetUserInfo).Methods("GET")


}

