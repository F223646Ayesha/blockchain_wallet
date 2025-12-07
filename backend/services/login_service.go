package services

import (
	"context"
	"crypto-wallet/config"
	"crypto-wallet/utils"
	"errors"
)

// -----------------------------------------------------
// 🚀 LOGIN SERVICE — Using Firestore OTP Verification
// -----------------------------------------------------
func LoginService(email, otp string) (string, map[string]interface{}, map[string]interface{}, error) {

	// STEP 1 — Verify OTP (new Firestore-based function)
	if !VerifyOTPService(email, otp) {
		return "", nil, nil, errors.New("invalid or expired OTP")
	}

	// STEP 2 — Fetch user from Firestore
	ctx := context.Background()

	userDoc, err := config.Firestore.Collection("users").
		Where("email", "==", email).
		Limit(1).
		Documents(ctx).Next()

	if err != nil {
		return "", nil, nil, errors.New("user not found")
	}

	user := userDoc.Data()
	walletID := user["wallet_id"].(string)

	// STEP 3 — Fetch wallet record
	walletDoc, err := config.Firestore.Collection("wallets").Doc(walletID).Get(ctx)
	if err != nil {
		return "", nil, nil, errors.New("wallet not found")
	}

	wallet := walletDoc.Data()

	// STEP 4 — Generate JWT token for session
	token := utils.GenerateToken(email)

	return token, user, wallet, nil
}
