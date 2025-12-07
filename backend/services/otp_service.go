package services

import (
	"context"
	"crypto-wallet/config"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
)

const otpExpiry = 5 * time.Minute

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendOTPService(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}

	ctx := context.Background()

	otp := generateOTP()

	// Save OTP to Firestore
	_, _, err := config.Firestore.Collection("otps").Add(ctx, map[string]interface{}{
		"email":     email,
		"otp":       otp,
		"createdAt": time.Now(),
	})
	if err != nil {
		return errors.New("failed to store OTP in Firestore")
	}

	subject := "Your Wallet Login OTP"
	body := fmt.Sprintf(
		"Your one-time login code is: %s\n\nThis code expires in 5 minutes.",
		otp,
	)

	// Send real email
	err = SendEmail(email, subject, body)
	if err != nil {
		return errors.New("failed to send email")
	}

	log.Println("OTP sent successfully to:", email)
	return nil
}

// -----------------------------------------------------
// VERIFY OTP (Correct Version — Firestore based)
// -----------------------------------------------------
func VerifyOTPService(email, otp string) bool {
	if email == "" || otp == "" {
		return false
	}

	ctx := context.Background()

	// Get latest OTP for this email
	q := config.Firestore.Collection("otps").
		Where("email", "==", email).
		OrderBy("createdAt", firestore.Desc).
		Limit(1).
		Documents(ctx)

	doc, err := q.Next()
	if err != nil {
		log.Println("No OTP found:", err)
		return false
	}

	data := doc.Data()
	storedOtp := data["otp"].(string)

	// Parse timestamp
	createdAt := data["createdAt"].(time.Time)

	// Check expiry
	if time.Since(createdAt) > otpExpiry {
		log.Println("OTP expired")
		return false
	}

	// Compare OTP
	return storedOtp == otp
}
