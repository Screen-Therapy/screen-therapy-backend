package models

// FriendRequest struct (no logic yet)
type FriendRequest struct {
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
	Status     string `json:"status"`
}
