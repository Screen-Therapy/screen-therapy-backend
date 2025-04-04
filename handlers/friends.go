package handlers

import "net/http"

// GetFriends placeholder function
func GetFriends(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Friends list feature coming soon!"))
}
