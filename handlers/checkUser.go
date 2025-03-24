package handlers

import (
	"context"
	"fmt"
	"net/http"

	"screen-therapy-backend/config"

	"github.com/gorilla/mux"
)

// CheckUser verifies if a user exists in Firestore
func CheckUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	// 🔍 Check if the user exists in Firestore
	doc, err := config.Client.Collection("users").Doc(userId).Get(context.Background())
	if err != nil || !doc.Exists() {
		w.WriteHeader(http.StatusNotFound) // 404 if user does not exist
		fmt.Fprintf(w, "User not found")
		return
	}

	// ✅ User exists
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User exists")
}
