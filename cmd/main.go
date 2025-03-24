package main

import (
	"fmt"

	"screen-therapy-backend/config"
)

func main() {
	// Initialize Firebase & Firestore
	config.InitFirebase()
	defer config.Client.Close()

	fmt.Println("🚀 Firestore connection established!")
}
