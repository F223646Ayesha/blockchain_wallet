package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var Firestore *firestore.Client

// ⭐ Required for OTP API
var FirebaseAPIKey string

func InitFirestore() {
	ctx := context.Background()

	// ===============================
	// 🔥 Load ENV variables
	// ===============================
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	FirebaseAPIKey = os.Getenv("FIREBASE_API_KEY")

	// ===============================
	// 🔥 Fixed path for Render
	// Render ALWAYS stores the uploaded secret file here:
	// /etc/secrets/<filename>
	// ===============================
	credentialsPath := "/etc/secrets/serviceAccountKey.json"

	// ===============================
	// 🔥 Validate configuration
	// ===============================
	if projectID == "" {
		log.Fatal("❌ FIREBASE_PROJECT_ID is missing")
	}

	if FirebaseAPIKey == "" {
		log.Fatal("❌ FIREBASE_API_KEY is missing")
	}

	// 🔥 We DO NOT depend on GOOGLE_APPLICATION_CREDENTIALS anymore
	// Render mounts secrets automatically in /etc/secrets/

	// ===============================
	// 🔥 Connect to Firestore
	// ===============================
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("❌ Firestore connection failed: %v", err)
	}

	Firestore = client
	log.Println("🔥 Firestore connected successfully")
}
