package services

import (
	"context"
	"crypto-wallet/config"
	"crypto-wallet/crypto"
	"crypto-wallet/database"
	"crypto-wallet/models"
	"crypto-wallet/utils"
	"encoding/hex"
	"errors"
	"fmt"
)

type TransactionInput struct {
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Amount    float64 `json:"amount"`
	Note      string  `json:"note"`
	Timestamp int64   `json:"timestamp"`
	Signature string  `json:"signature"`
}

// ============================================================
// Build canonical signing payload (MUST MATCH FRONTEND)
// ============================================================
func buildPayload(t TransactionInput) string {
	payload := t.Sender + "|" +
		t.Receiver + "|" +
		utils.FloatToString(t.Amount) + "|" +
		utils.IntToString(t.Timestamp) + "|" +
		t.Note

	fmt.Println("🟣 BACKEND buildPayload():", payload)

	return payload
}

// ============================================================
// Verify Digital Signature — logs added
// ============================================================
func VerifyTransactionSignature(tx TransactionInput, pubKeyHex string) error {

	fmt.Println("======================================================")
	fmt.Println("🔵 VerifyTransactionSignature CALLED")
	fmt.Println("🔵 Frontend Signature (hex):", tx.Signature)
	fmt.Println("🔵 Public Key (hex):", pubKeyHex)

	pubKey, err := crypto.PublicKeyFromHex(pubKeyHex)
	if err != nil {
		fmt.Println("❌ ERROR parsing public key:", err)
		return err
	}

	signBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		fmt.Println("❌ ERROR decoding signature hex:", err)
		return err
	}

	fmt.Println("🔵 Signature (decoded bytes):", signBytes)

	backendPayload := buildPayload(tx)
	fmt.Println("🔵 Payload used for verification:", backendPayload)

	ok := crypto.VerifySignature(*pubKey, backendPayload, signBytes)
	fmt.Println("🔵 Signature Verification Result:", ok)
	fmt.Println("======================================================")

	if !ok {
		return utils.ErrInvalidSignature
	}

	return nil
}

// ============================================================
// Main Transaction Flow — logs added
// ============================================================
func SendTransactionService(tx TransactionInput) error {

	fmt.Println("======================================================")
	fmt.Println("🚀 SendTransactionService CALLED")
	fmt.Printf("🚀 TX Received: %+v\n", tx)

	ctx := context.Background()

	// 1. Fetch sender wallet
	doc, err := config.Firestore.Collection("wallets").
		Where("wallet_id", "==", tx.Sender).
		Limit(1).
		Documents(ctx).
		Next()

	if err != nil {
		fmt.Println("❌ Sender wallet NOT FOUND")
		return errors.New("sender wallet not found")
	}

	pubKeyHex := doc.Data()["public_key"].(string)
	fmt.Println("🔵 Sender Public Key (from Firestore):", pubKeyHex)

	// 2. Verify signature
	err = VerifyTransactionSignature(tx, pubKeyHex)
	if err != nil {
		fmt.Println("❌ DIGITAL SIGNATURE FAILED:", err)
		return err
	}
	fmt.Println("✅ DIGITAL SIGNATURE VERIFIED!")

	// 3. Select UTXOs
	selectedUTXOs, changeAmount, err := SelectUTXOs(tx.Sender, tx.Amount)
	if err != nil {
		fmt.Println("❌ ERROR Selecting UTXOs:", err)
		return err
	}

	fmt.Println("🔵 Selected UTXOs:", selectedUTXOs)
	fmt.Println("🔵 Change Amount:", changeAmount)

	// 4. Mark UTXOs as spent
	for _, u := range selectedUTXOs {
		_ = database.MarkUTXOSpent(u.ID)
	}

	// 5. Create receiver UTXO
	receiverUTXO := models.UTXO{
		WalletID: tx.Receiver,
		Amount:   tx.Amount,
		Spent:    false,
		TxID:     "pending",
	}
	database.AddUTXO(receiverUTXO)

	// 6. Handle change UTXO
	if changeAmount > 0 {
		changeUTXO := models.UTXO{
			WalletID: tx.Sender,
			Amount:   changeAmount,
			Spent:    false,
			TxID:     "pending",
		}
		database.AddUTXO(changeUTXO)
	}

	// 7. Save pending transaction
	_, _, err = config.Firestore.Collection("pending_transactions").Add(ctx, map[string]interface{}{
		"sender":    tx.Sender,
		"receiver":  tx.Receiver,
		"amount":    tx.Amount,
		"note":      tx.Note,
		"timestamp": tx.Timestamp,
		"signature": tx.Signature,
	})

	fmt.Println("✅ Transaction saved in pending_transactions")
	fmt.Println("======================================================")

	return err
}
