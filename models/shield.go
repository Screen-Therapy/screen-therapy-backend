package models

import "time"

type BlockType string

const (
	BlockTypeStrict     BlockType = "strict"
	BlockTypeSemiStrict BlockType = "semi_strict"
	BlockTypeNotStrict  BlockType = "not_strict"
)

type ChallengeType string

const (
	ChallengeNone    ChallengeType = "none"
	ChallengeMath    ChallengeType = "math"
	ChallengeTyping  ChallengeType = "typing"
	ChallengeCustom  ChallengeType = "custom"
)

type Shield struct {
	ID               string         `json:"id" firestore:"id"`
	OwnerUserID      string         `json:"ownerUserId" firestore:"ownerUserId"`
	CreatedByUserID  string         `json:"createdByUserId" firestore:"createdByUserId"`
	AppGroupName     string         `json:"appGroupName" firestore:"appGroupName"`
	BlockedApps      []string       `json:"blockedApps" firestore:"blockedApps"` // app bundle IDs or domains
	BlockType        BlockType      `json:"blockType" firestore:"blockType"` // strict, semi, etc.
	ChallengeType    ChallengeType  `json:"challengeType" firestore:"challengeType"`
	TotalDurationMin int            `json:"totalDurationMin" firestore:"totalDurationMin"` // in minutes
	StartTime        *time.Time     `json:"startTime,omitempty" firestore:"startTime,omitempty"`
	EndTime          *time.Time     `json:"endTime,omitempty" firestore:"endTime,omitempty"`
	RepeatDays       []string       `json:"repeatDays" firestore:"repeatDays"` // e.g., ["monday", "tuesday"]
	Quote            string         `json:"quote,omitempty" firestore:"quote"`
	Icon             string         `json:"icon" firestore:"icon"`           // SF Symbol name
	PrimaryColor     string         `json:"primaryColor" firestore:"primaryColor"`   // hex string
	SecondaryColor   string         `json:"secondaryColor" firestore:"secondaryColor"` // hex string
	Reason           string         `json:"reason,omitempty" firestore:"reason"`     // optional
	Status           string         `json:"status" firestore:"status"` // pending, active, expired
	CreatedAt        time.Time      `json:"createdAt" firestore:"createdAt"`
}