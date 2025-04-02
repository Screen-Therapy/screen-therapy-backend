package friends

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"screen-therapy-backend/config"

	"cloud.google.com/go/firestore"
)

type AddFriendRequest struct {
	UserID     string `json:"userId"`
	FriendCode string `json:"friendCode"`
}

// AddFriendHandler lets a user add another by friendCode
func AddFriendHandler(w http.ResponseWriter, r *http.Request) {
	var req AddFriendRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.FriendCode == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	client := config.Client

	// Find friend by friendCode
	friendSnap, err := client.Collection("users").Where("friendCode", "==", req.FriendCode).Limit(1).Documents(ctx).Next()
	if err != nil {
		http.Error(w, "Friend not found", http.StatusNotFound)
		return
	}
	friendID := friendSnap.Ref.ID

	// Confirm both users exist
	userRef := client.Collection("users").Doc(req.UserID)
	friendRef := client.Collection("users").Doc(friendID)

	_, err = userRef.Get(ctx)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Add friend to user's friends subcollection
	_, err = userRef.Collection("friends").Doc(friendID).Set(ctx, map[string]interface{}{
		"userId": friendID,
		"addedAt": firestore.ServerTimestamp,
	})
	if err != nil {
		http.Error(w, "Failed to add friend", http.StatusInternalServerError)
		return
	}

	// Add user to friend's friends subcollection
	_, err = friendRef.Collection("friends").Doc(req.UserID).Set(ctx, map[string]interface{}{
		"userId": req.UserID,
		"addedAt": firestore.ServerTimestamp,
	})
	if err != nil {
		http.Error(w, "Failed to update friend's list", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "✅ Friends added successfully")
}

// GetFriendsHandler returns a list of friends for the given userId
func GetFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "Missing userId", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	friendsSnap, err := config.Client.Collection("users").Doc(userId).Collection("friends").Documents(ctx).GetAll()
	if err != nil {
		http.Error(w, "Failed to retrieve friends list", http.StatusInternalServerError)
		return
	}

	friends := []map[string]interface{}{}
	for _, doc := range friendsSnap {
		friends = append(friends, doc.Data())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}
