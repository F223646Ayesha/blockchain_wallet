package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"encoding/json"
)

type TransactionInput struct {
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Amount    float64 `json:"amount"`
	Note      string  `json:"note"`
	Timestamp int64   `json:"timestamp"`
}

func buildPayload(tx TransactionInput) string {
	return fmt.Sprintf("%s|%s|%f|%d|%s",
		tx.Sender,
		tx.Receiver,
		tx.Amount,
		tx.Timestamp,
		tx.Note,
	)
}

func main() {

	// === 1. Define transaction ===
	tx := TransactionInput{
		Sender:    "YOUR_SENDER_WALLET",
		Receiver:  "YOUR_RECEIVER_WALLET",
		Amount:    10,
		Note:      "testing",
		Timestamp: 1736010000,
	}

	fmt.Println("Payload:")
	fmt.Println(buildPayload(tx))

	// === 2. Load real private key (hex from Firestore) ===
	privHex := "YOUR_PRIVATE_KEY_HEX"

	privBytes, _ := hex.DecodeString(privHex)
	d := new(big.Int).SetBytes(privBytes)

	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = elliptic.P256()
	privateKey.D = d
	privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(privBytes)

	// === 3. Hash payload ===
	payload := buildPayload(tx)
	hash := sha256.Sum256([]byte(payload))

	// === 4. Sign payload ===
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		log.Fatal(err)
	}

	// Combine r + s → hex signature
	sig := append(r.Bytes(), s.Bytes()...)
	sigHex := hex.EncodeToString(sig)

	fmt.Println("\nSignature (hex):")
	fmt.Println(sigHex)

	// === 5. Print JSON to send to API ===
	output := map[string]interface{}{
		"sender":    tx.Sender,
		"receiver":  tx.Receiver,
		"amount":    tx.Amount,
		"note":      tx.Note,
		"timestamp": tx.Timestamp,
		"signature": sigHex,
	}

	jsonOutput, _ := json.MarshalIndent(output, "", "  ")
	fmt.Println("\n\nSend this JSON to POST /api/transaction/send:")
	fmt.Println(string(jsonOutput))
}
