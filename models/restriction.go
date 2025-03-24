package models

// Restriction struct (no logic yet)
type Restriction struct {
	UserID string `json:"userId"`
	App    string `json:"app"`
	Limit  int    `json:"limit"`
}
