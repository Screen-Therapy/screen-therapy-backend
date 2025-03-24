package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

var client *firestore.Client

func initFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase: %v", err)
	}

	client, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("Error initializing Firestore: %v", err)
	}
	fmt.Println("✅ Firestore initialized successfully!")
}

// FriendRequest struct
type FriendRequest struct {
	SenderID     string         `json:"senderId"`
	ReceiverID   string         `json:"receiverId"`
	Restrictions map[string]int `json:"restrictions"`
}

// API Handler to send friend request
func sendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req FriendRequest

	// Decode JSON request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Store in Firestore
	_, _, err = client.Collection("requests").Add(context.Background(), map[string]interface{}{
		"sender":      req.SenderID,
		"receiver":    req.ReceiverID,
		"restrictions": req.Restrictions,
		"status":      "pending",
	})
	if err != nil {
		log.Printf("❌ Firestore error: %v", err) // Log the error for debugging
		http.Error(w, "Error saving to Firestore", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Friend request sent successfully!")
}


// Main function to start the server
func main() {
	initFirebase()
	defer client.Close()

	// Setup router
	r := mux.NewRouter()
	r.HandleFunc("/sendRequest", sendRequestHandler).Methods("POST")

	fmt.Println("🚀 Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
