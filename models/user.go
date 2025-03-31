package models

type User struct {
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username,omitempty"` // Only if set later
}
