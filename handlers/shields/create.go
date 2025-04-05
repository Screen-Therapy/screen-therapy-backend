package shields

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"screen-therapy-backend/config"

	"github.com/google/uuid"
)

// CreateShieldHandler handles POST /shields
func CreateShieldHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	fmt.Println("📥 Incoming request to create shield")

	var shield Shield
	if err := json.NewDecoder(r.Body).Decode(&shield); err != nil {
		fmt.Println("❌ Failed to decode request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	fmt.Printf("🛡️ Shield decoded: %+v\n", shield)

	// 🔍 Basic validation
	if shield.OwnerUserID == ""|| shield.AppGroupName == "" {
		fmt.Println("⚠️ Missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// 🔐 Validate block type
	switch shield.BlockType {
	case BlockTypeStrict, BlockTypeSemiStrict, BlockTypeNotStrict:
		fmt.Println("✅ Block type is valid:", shield.BlockType)
	default:
		fmt.Println("❌ Invalid block type:", shield.BlockType)
		http.Error(w, "Invalid block type", http.StatusBadRequest)
		return
	}

	// 🧠 Validate challenge type
	switch shield.ChallengeType {
	case ChallengeNone, ChallengeMath, ChallengeTyping, ChallengeCustom:
		fmt.Println("✅ Challenge type is valid:", shield.ChallengeType)
	default:
		fmt.Println("❌ Invalid challenge type:", shield.ChallengeType)
		http.Error(w, "Invalid challenge type", http.StatusBadRequest)
		return
	}

	// 🆔 ID + createdAt
	shield.ID = uuid.NewString()
	shield.CreatedAt = time.Now().UTC()
	fmt.Println("🆔 Assigned new ID and timestamp:", shield.ID, shield.CreatedAt)

	// Fallback to self-created
	if shield.CreatedByUserID == "" {
		shield.CreatedByUserID = shield.OwnerUserID
		fmt.Println("🔄 CreatedByUserID fallback to OwnerUserID:", shield.CreatedByUserID)
	}

	// Default status
	if shield.Status == "" {
		shield.Status = "pending"
		fmt.Println("🕒 Defaulting shield status to 'pending'")
	}

	// 💾 Save to Firestore
	_, err := config.Client.Collection("shields").Doc(shield.ID).Set(ctx, shield)
	if err != nil {
		fmt.Println("❌ Firestore error:", err)
		http.Error(w, "Failed to create shield", http.StatusInternalServerError)
		return
	}
	fmt.Println("✅ Shield successfully stored in Firestore:", shield.ID)

	// ✅ Return created shield
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(shield); err != nil {
		fmt.Println("❌ Failed to encode response:", err)
	}
}
