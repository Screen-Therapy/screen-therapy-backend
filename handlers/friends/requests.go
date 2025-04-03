package friends

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"screen-therapy-backend/config"

	"cloud.google.com/go/firestore"
)

type AddFriendRequest struct {
	UserID     string `json:"userId"`
	FriendCode string `json:"friendCode"`
}

type Friend struct {
	UserID   string `json:"userId"`  // ✅ must be "userId" to match Swift
	Username string `json:"username"`
}


// AddFriendHandler lets a user add another by friendCode
func AddFriendHandler(w http.ResponseWriter, r *http.Request) {
	var req AddFriendRequest

	log.Println("📥 Incoming friend request")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.UserID == "" || req.FriendCode == "" {
		log.Printf("❌ Invalid input: %+v\n", req)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Looking up friendCode: %s\n", req.FriendCode)
	ctx := context.Background()
	client := config.Client

	// Find friend by friendCode
	friendSnap, err := client.Collection("users").Where("friendCode", "==", req.FriendCode).Limit(1).Documents(ctx).Next()
	if err != nil {
		log.Println("❌ Friend not found")
		http.Error(w, "Friend not found", http.StatusNotFound)
		return
	}
	friendID := friendSnap.Ref.ID
	log.Printf("✅ Friend found: %s\n", friendID)

	// Confirm both users exist
	userRef := client.Collection("users").Doc(req.UserID)
	userSnap, err := userRef.Get(ctx)
	if err != nil {
		log.Printf("❌ User not found: %s\n", req.UserID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	friendRef := client.Collection("users").Doc(friendID)

	// Add friend to user's friends subcollection
	log.Println("🔗 Adding friend to user's friends list...")
	_, err = userRef.Collection("friends").Doc(friendID).Set(ctx, map[string]interface{}{
		"userId":   friendID,
		"username": friendSnap.Data()["username"],
		"addedAt":  firestore.ServerTimestamp,
	})
	if err != nil {
		log.Printf("❌ Failed to add friend: %v\n", err)
		http.Error(w, "Failed to add friend", http.StatusInternalServerError)
		return
	}

	// Add user to friend's friends subcollection
	log.Println("🔁 Adding user to friend's friends list...")
	_, err = friendRef.Collection("friends").Doc(req.UserID).Set(ctx, map[string]interface{}{
		"userId":   req.UserID,
		"username": userSnap.Data()["username"],
		"addedAt":  firestore.ServerTimestamp,
	})
	if err != nil {
		log.Printf("❌ Failed to update friend's list: %v\n", err)
		http.Error(w, "Failed to update friend's list", http.StatusInternalServerError)
		return
	}

	log.Println("✅ Friends added successfully")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Friends added successfully",
	})
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

	var friends []Friend
	for _, doc := range friendsSnap {
		data := doc.Data()
		friend := Friend{
			UserID:   fmt.Sprintf("%v", data["userId"]), // ✅ make sure this is included
			Username: fmt.Sprintf("%v", data["username"]),
		}
		friends = append(friends, friend)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(friends)
}
