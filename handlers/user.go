package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"screen-therapy-backend/config"

	"github.com/gorilla/mux"
)

// GetUserInfo fetches the username and friendCode for a given userId
func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	fmt.Printf("🔍 Fetching user info for userId: %s\n", userId)

	doc, err := config.Client.Collection("users").Doc(userId).Get(context.Background())
	if err != nil {
		fmt.Printf("❌ Firestore error: %v\n", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !doc.Exists() {
		fmt.Println("⚠️ User document does not exist.")
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data := doc.Data()
	fmt.Printf("📦 Retrieved Firestore data: %+v\n", data)

	username, ok1 := data["username"]
	friendCode, ok2 := data["friendCode"]

	if !ok1 || !ok2 {
		fmt.Println("❗ Missing expected fields in document.")
	}

	response := map[string]interface{}{
		"username":   username,
		"friendCode": friendCode,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Printf("❌ JSON encode error: %v\n", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Successfully sent user info response.")
}
