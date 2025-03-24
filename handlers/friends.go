package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"screen-therapy-backend/config"
)

// Friend structure
type Friend struct {
	UserID   string `json:"userId"`
	FriendID string `json:"friendId"`
}

// GetFriends - Retrieve a list of a user's friends
func GetFriends(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")

	if userID == "" {
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	// Query Firestore for all friends of this user
	docs, err := config.Client.Collection("friends").Where("userId", "==", userID).Documents(context.Background()).GetAll()
	if err != nil {
		log.Printf("Firestore error: %v", err)
		http.Error(w, "Error retrieving friends", http.StatusInternalServerError)
		return
	}

	var friends []Friend
	for _, doc := range docs {
		var friend Friend
		doc.DataTo(&friend)
		friends = append(friends, friend)
	}

	// Return friends as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}

// AddFriend - Add a new friend
func AddFriend(w http.ResponseWriter, r *http.Request) {
	var req Friend

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Store friendship in Firestore
	_, _, err = config.Client.Collection("friends").Add(context.Background(), req)
	if err != nil {
		log.Printf("Firestore error: %v", err)
		http.Error(w, "Error saving friend", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Friend added successfully!")
}
