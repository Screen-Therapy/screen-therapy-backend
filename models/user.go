package models

import "time"

type User struct {
	UserID       string                 `json:"userId"`
	Email        string                 `json:"email"`
	Username     string                 `json:"username,omitempty"`
	FriendCode   string                 `json:"friendCode,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	Roles        map[string]interface{} `json:"roles,omitempty"`
	AccountableTo []string              `json:"accountableTo,omitempty"`
	GuardianOf   []string               `json:"guardianOf,omitempty"`
}
