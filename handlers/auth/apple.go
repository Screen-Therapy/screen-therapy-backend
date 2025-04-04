package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"screen-therapy-backend/config"
	"screen-therapy-backend/models"

	"cloud.google.com/go/firestore"

	"github.com/gorilla/mux"
)

// RegisterAppleUser registers a new Apple user if they don't already exist
func RegisterAppleUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	doc, err := config.Client.Collection("users").Doc(user.UserID).Get(context.Background())
	if err == nil && doc.Exists() {
		fmt.Fprintf(w, "User already exists")
		return
	}

	// Generate friend code for invites
	friendCode := generateFriendCode()


	_, err = config.Client.Collection("users").Doc(user.UserID).Set(context.Background(), map[string]interface{}{
		"userId":       user.UserID,
		"email":        user.Email,
		"friendCode":   friendCode,
		"roles": map[string]interface{}{
			"guardianOf":    []string{},
			"accountableTo": []string{},
		},
		"createdAt": firestore.ServerTimestamp,
	})
	
	if err != nil {
		log.Printf("❌ Firestore error: %v", err)
		http.Error(w, "Failed to store Apple user", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "✅ Apple user registered successfully")
}

// CheckAppleUser verifies if an Apple user exists in Firestore
func CheckAppleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	doc, err := config.Client.Collection("users").Doc(userId).Get(context.Background())
	if err != nil || !doc.Exists() {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Apple user not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Apple user exists")
}

// CheckAppleUsername checks if the Apple user has a username set
func CheckAppleUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	doc, err := config.Client.Collection("users").Doc(userId).Get(context.Background())
	if err != nil || !doc.Exists() {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User not found")
		return
	}

	username, ok := doc.Data()["username"].(string)
	if ok && username != "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Username exists")
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Username not set")
	}
}

// SetAppleUsername adds or updates the username field for an Apple user
func SetAppleUsername(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID   string `json:"userId"`
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || data.UserID == "" || data.Username == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	_, err = config.Client.Collection("users").Doc(data.UserID).Update(context.Background(), []firestore.Update{
		{Path: "username", Value: data.Username},
	})
	
	if err != nil {
		log.Printf("Error updating username: %v", err)
		http.Error(w, "Failed to set username", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "✅ Username set for Apple user")
}


func generateFriendCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
