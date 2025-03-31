package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"screen-therapy-backend/config"
	"screen-therapy-backend/models"
)

// RegisterEmailUser saves a new email user in Firestore (called after Firebase Auth sign-up)
func RegisterEmailUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	log.Println("📥 Incoming request to register email user")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil || user.UserID == "" || user.Email == "" || user.Username == "" {
		log.Printf("❌ Invalid input: %+v\n", user)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Checking if user already exists: %s\n", user.UserID)
	doc, err := config.Client.Collection("users").Doc(user.UserID).Get(context.Background())
	if err == nil && doc.Exists() {
		log.Println("⚠️ User already exists")
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	log.Printf("📝 Registering user: %+v\n", user)
	_, err = config.Client.Collection("users").Doc(user.UserID).Set(context.Background(), map[string]interface{}{
		"userId":   user.UserID,
		"email":    user.Email,
		"username": user.Username,
	})

	if err != nil {
		log.Printf("❌ Firestore error: %v\n", err)
		http.Error(w, "Failed to store user", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Email user registered successfully")
	fmt.Fprintf(w, "✅ Email user registered successfully")
}


// LoginEmailUser checks if the user exists (called after Firebase Auth login)
func LoginEmailUser(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID string `json:"userId"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || data.UserID == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	doc, err := config.Client.Collection("users").Doc(data.UserID).Get(context.Background())
	if err != nil || !doc.Exists() {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "✅ Email user exists")
}
