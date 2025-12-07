package services

import (
	"context"
	"crypto-wallet/config"
	"crypto-wallet/utils"
	"errors"
	"time"
)

type RegisterResponse struct {
	Token  string                 `json:"token"`
	User   map[string]interface{} `json:"user"`
	Wallet map[string]interface{} `json:"wallet"`
}

func RegisterUserService(name, email, cnic string) (*RegisterResponse, error) {

	ctx := context.Background()

	// 1. Check if user already exists
	existing, _ := config.Firestore.Collection("users").
		Where("email", "==", email).
		Limit(1).
		Documents(ctx).Next()

	if existing != nil {
		return nil, errors.New("email already registered")
	}

	// 2. Create new wallet
	wallet, err := CreateWalletService()
	if err != nil {
		return nil, err
	}

	// 3. Store user WITHOUT public/private key copies
	userData := map[string]interface{}{
		"name":       name,
		"email":      email,
		"cnic":       cnic,
		"wallet_id":  wallet.WalletID,
		"created_at": time.Now().Unix(),
	}

	_, _, err = config.Firestore.Collection("users").Add(ctx, userData)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// 4. Generate JWT token
	token := utils.GenerateToken(email)

	// 5. Wallet payload for frontend (matches your JS: wallet.wallet_id, wallet.public_key, wallet.private_key_enc)
	walletPayload := map[string]interface{}{
		"wallet_id":       wallet.WalletID,
		"public_key":      wallet.PublicKey,  // already HEX string
		"private_key_enc": wallet.PrivateKey, // encrypted PKCS8 string
	}

	// 6. Final response
	return &RegisterResponse{
		Token:  token,
		User:   userData,
		Wallet: walletPayload,
	}, nil
}
