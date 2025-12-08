package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var Firestore *firestore.Client

// Firebase Web API key for sending OTP
var FirebaseAPIKey string

func InitFirestore() {
	ctx := context.Background()

	// Read environment variables
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	credentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	FirebaseAPIKey = os.Getenv("FIREBASE_API_KEY")

	// =============================
	// DEBUG LOGS (REQUIRED FOR RENDER)
	// =============================
	log.Println("🔍 FIRESTORE INIT STARTING...")
	log.Println("📌 FIREBASE_PROJECT_ID =", projectID)
	log.Println("📌 FIREBASE_API_KEY =", FirebaseAPIKey)
	log.Println("📌 GOOGLE_APPLICATION_CREDENTIALS =", credentials)

	// ---------- Check credentials path ----------
	if credentials == "" {
		log.Println("❌ ERROR: GOOGLE_APPLICATION_CREDENTIALS is EMPTY")
	} else {
		// Check if file exists at given path
		if _, err := os.Stat(credentials); os.IsNotExist(err) {
			log.Println("❌ ERROR: Credentials file NOT FOUND at:", credentials)
			log.Println("💡 TIP: On Render, file must be named EXACTLY 'serviceAccountKey.json'")
			log.Println("💡 TIP: And GOOGLE_APPLICATION_CREDENTIALS must be '/etc/secrets/serviceAccountKey.json'")
		} else {
			log.Println("✔ Credentials file FOUND at:", credentials)
		}
	}

	// ---------- Validate required variables ----------
	if projectID == "" {
		log.Fatal("❌ FIREBASE_PROJECT_ID is missing")
	}
	if FirebaseAPIKey == "" {
		log.Fatal("❌ FIREBASE_API_KEY is missing")
	}
	if credentials == "" {
		log.Fatal("❌ GOOGLE_APPLICATION_CREDENTIALS is missing")
	}

	// ---------- Initialize Firestore ----------
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentials))
	if err != nil {
		log.Fatalf("❌ Firestore connection failed: %v", err)
	}

	Firestore = client
	log.Println("🔥 Firestore connected successfully")
}
