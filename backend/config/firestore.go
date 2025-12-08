func InitFirestore() {
	ctx := context.Background()

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	FirebaseAPIKey = os.Getenv("FIREBASE_API_KEY")

	// 🔥 Use fixed Render secret file path
	credentialsPath := "/etc/secrets/serviceAccountKey.json"

	if projectID == "" || FirebaseAPIKey == "" {
		log.Fatal("❌ Missing FIREBASE_PROJECT_ID or FIREBASE_API_KEY")
	}

	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		log.Fatalf("❌ Firestore connection failed: %v", err)
	}

	Firestore = client
	log.Println("🔥 Firestore connected successfully")
}
