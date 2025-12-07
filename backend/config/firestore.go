package config

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

var Firestore *firestore.Client

// ⭐ Add Firebase API Key for OTP login
var FirebaseAPIKey string

func InitFirestore() {
	ctx := context.Background()

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	credentials := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	FirebaseAPIKey = os.Getenv("FIREBASE_API_KEY") // ⭐ NEW

	if projectID == "" || credentials == "" || FirebaseAPIKey == "" {
		log.Fatal("❌ Missing FIREBASE_PROJECT_ID / GOOGLE_APPLICATION_CREDENTIALS / FIREBASE_API_KEY")
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentials))
	if err != nil {
		log.Fatalf("❌ Firestore connection failed: %v", err)
	}

	Firestore = client
	log.Println("🔥 Firestore connected successfully")
}
