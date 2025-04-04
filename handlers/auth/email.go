package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"screen-therapy-backend/config"
	"screen-therapy-backend/models"

	"cloud.google.com/go/firestore"
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

	doc, err := config.Client.Collection("users").Doc(user.UserID).Get(context.Background())
	if err == nil && doc.Exists() {
		log.Println("⚠️ User already exists")
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	friendCode := generateFriendCode()

	_, err = config.Client.Collection("users").Doc(user.UserID).Set(context.Background(), map[string]interface{}{
		"userId":     user.UserID,
		"email":      user.Email,
		"username":   user.Username,
		"friendCode": friendCode,
		"roles": map[string]interface{}{
			"guardianOf":    []string{},
			"accountableTo": []string{},
		},
		"createdAt": firestore.ServerTimestamp,
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

