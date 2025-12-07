package utils

import (
	"context"
	"crypto-wallet/config"
	"time"

	"cloud.google.com/go/firestore"
)

func CheckOTP(email, otp string) bool {
	ctx := context.Background()

	iter := config.Firestore.Collection("otps").
		Where("email", "==", email).
		OrderBy("createdAt", firestore.Desc).
		Limit(1).
		Documents(ctx)

	doc, err := iter.Next()
	if err != nil {
		return false
	}

	data := doc.Data()
	savedOtp := data["otp"].(string)

	// --- FIX: robust casting for Firestore number ---
	var savedAt int64
	switch v := data["createdAt"].(type) {
	case int64:
		savedAt = v
	case int:
		savedAt = int64(v)
	case float64:
		savedAt = int64(v)
	default:
		return false
	}

	// Expiry check (5 minutes)
	if time.Now().Unix()-savedAt > 300 {
		return false
	}

	return otp == savedOtp
}
