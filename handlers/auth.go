package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"screen-therapy-backend/config"
)

// User struct for storing user data
type User struct {
	UserID   string `json:"userId"`
	FullName string `json:"fullName"`
	Email    string `json:"email"`
}

// RegisterUser saves a new user in Firestore if they don’t already exist
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode incoming JSON request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// 🔍 Check if user already exists
	doc, err := config.Client.Collection("users").Doc(user.UserID).Get(context.Background())
	if err == nil && doc.Exists() {
		fmt.Fprintf(w, "User already exists")
		return
	}

	// 📝 Save the user in Firestore
	_, err = config.Client.Collection("users").Doc(user.UserID).Set(context.Background(), map[string]interface{}{
		"userId":   user.UserID,
		"fullName": user.FullName,
		"email":    user.Email,
	})

	if err != nil {
		log.Printf("❌ Firestore error: %v", err)
		http.Error(w, "Failed to store user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "✅ User registered successfully")
}
