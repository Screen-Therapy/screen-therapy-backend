package config

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

// Global Firestore client
var Client *firestore.Client

// InitFirebase connects to Firestore
func InitFirebase() {
	opt := option.WithCredentialsFile("serviceAccountKey.json")

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("❌ Error initializing Firebase: %v", err)
	}

	Client, err = app.Firestore(context.Background())
	if err != nil {
		log.Fatalf("❌ Error connecting to Firestore: %v", err)
	}

	log.Println("✅ Connected to Firestore!")
}
