package services

import (
	"context"
	"crypto-wallet/config"
	"crypto-wallet/crypto"
	"crypto-wallet/database" // ⬅ REQUIRED
	"crypto-wallet/models"   // ⬅ REQUIRED
	"crypto-wallet/utils"
	"encoding/base64"
	"encoding/hex"
	"errors"

	"cloud.google.com/go/firestore"
)

var AES_SECRET = []byte("12345678901234567890123456789012") // 32-byte AES key

type CreateWalletResponse struct {
	WalletID   string `json:"wallet_id"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key_encrypted"`
}

// ====================================================
// CREATE WALLET (Correct PKCS8 + AES encryption)
// ====================================================
func CreateWalletService() (*CreateWalletResponse, error) {

	// 1️⃣ Generate ECDSA keypair (must include PKCS8 bytes)
	keypair, err := crypto.GenerateKeyPair()
	if err != nil {
		return nil, errors.New("keypair generation failed")
	}

	// 2️⃣ Convert public key (SEC1) → hex
	pubHex := hex.EncodeToString(keypair.PublicKey)

	// 3️⃣ Wallet ID = SHA256(public key)
	walletID := crypto.GenerateWalletID(keypair.PublicKey)

	// 4️⃣ Convert PKCS8 → Base64
	if len(keypair.PrivateKeyPKCS8) == 0 {
		return nil, errors.New("PKCS8 key not generated — FIX GenerateKeyPair()")
	}
	pkcs8B64 := base64.StdEncoding.EncodeToString(keypair.PrivateKeyPKCS8)

	// 5️⃣ Encrypt Base64 PKCS8 using AES-256
	encryptedKey, err := crypto.EncryptAES(AES_SECRET, pkcs8B64)
	if err != nil {
		return nil, err
	}

	// 6️⃣ Store wallet in Firestore
	ctx := context.Background()

	_, err = config.Firestore.Collection("wallets").
		Doc(walletID).
		Set(ctx, map[string]interface{}{
			"wallet_id":       walletID,
			"public_key":      pubHex,
			"private_key_enc": encryptedKey,
			"beneficiaries":   []string{},
		})

	if err != nil {
		return nil, err
	}

	// =======================================================
	// ⭐ 7️⃣ ADD GENESIS UTXO — Give wallet initial balance
	// =======================================================
	genesis := models.UTXO{
		ID:       "",
		WalletID: walletID,
		Amount:   100.0,
		Spent:    false,
		TxID:     "genesis",
	}

	_ = database.AddUTXO(genesis) // ignore error; UTXO DB is local

	// 8️⃣ Final response
	return &CreateWalletResponse{
		WalletID:   walletID,
		PublicKey:  pubHex,
		PrivateKey: encryptedKey,
	}, nil
}

// ====================================================
// GET WALLET
// ====================================================
func GetWalletService(walletID string) (map[string]interface{}, error) {

	ctx := context.Background()
	doc, err := config.Firestore.Collection("wallets").
		Doc(walletID).
		Get(ctx)

	if err != nil || !doc.Exists() {
		return nil, utils.ErrInvalidWallet
	}

	return doc.Data(), nil
}

// ====================================================
// ADD BENEFICIARY
// ====================================================
func AddBeneficiaryService(walletID, beneficiary string) error {

	ctx := context.Background()

	doc, err := config.Firestore.Collection("wallets").
		Doc(walletID).
		Get(ctx)

	if err != nil || !doc.Exists() {
		return errors.New("wallet not found")
	}

	data := doc.Data()
	arr, ok := data["beneficiaries"].([]interface{})
	if !ok {
		arr = []interface{}{}
	}

	// prevent duplicates
	for _, v := range arr {
		if v.(string) == beneficiary {
			return errors.New("beneficiary already exists")
		}
	}

	arr = append(arr, beneficiary)

	_, err = doc.Ref.Update(ctx, []firestore.Update{
		{Path: "beneficiaries", Value: arr},
	})

	return err
}

// ====================================================
// REMOVE BENEFICIARY
// ====================================================
func RemoveBeneficiaryService(walletID, beneficiary string) error {

	ctx := context.Background()

	doc, err := config.Firestore.Collection("wallets").
		Doc(walletID).
		Get(ctx)

	if err != nil || !doc.Exists() {
		return errors.New("wallet not found")
	}

	data := doc.Data()
	arr, ok := data["beneficiaries"].([]interface{})
	if !ok {
		return nil
	}

	newArr := []string{}
	for _, v := range arr {
		if v.(string) != beneficiary {
			newArr = append(newArr, v.(string))
		}
	}

	_, err = doc.Ref.Update(ctx, []firestore.Update{
		{Path: "beneficiaries", Value: newArr},
	})

	return err
}

// ====================================================
// UPDATE WALLET PROFILE (Name, Email, CNIC)
// ====================================================
func UpdateWalletProfileService(walletID, name, email, cnic string) error {

	ctx := context.Background()

	doc, err := config.Firestore.Collection("wallets").
		Doc(walletID).
		Get(ctx)

	if err != nil || !doc.Exists() {
		return errors.New("wallet not found")
	}

	updates := []firestore.Update{}

	if name != "" {
		updates = append(updates, firestore.Update{Path: "name", Value: name})
	}
	if email != "" {
		updates = append(updates, firestore.Update{Path: "email", Value: email})
	}
	if cnic != "" {
		updates = append(updates, firestore.Update{Path: "cnic", Value: cnic})
	}

	if len(updates) == 0 {
		return errors.New("nothing to update")
	}

	_, err = doc.Ref.Update(ctx, updates)
	return err
}
