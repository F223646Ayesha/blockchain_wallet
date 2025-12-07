package services

import (
	"context"
	"crypto-wallet/config"
	"crypto-wallet/database"
	"fmt"
	"time"
)

const ZAKAT_RATE = 0.025
const ZAKAT_POOL = "ZAKAT_POOL_001"

// ----------------------------------------------------------
// RUN MONTHLY ZAKAT — Creates Pending Transactions
// ----------------------------------------------------------
func RunZakatService() error {

	ctx := context.Background()

	// 1. Get all wallets
	iter := config.Firestore.Collection("wallets").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		walletID := doc.Data()["wallet_id"].(string)

		// Skip zakat pool itself
		if walletID == ZAKAT_POOL {
			continue
		}

		// 2. Calculate balance
		balance, err := GetBalance(walletID)
		if err != nil || balance <= 0 {
			continue
		}

		// 3. Compute zakat
		zakatAmount := balance * ZAKAT_RATE
		if zakatAmount <= 0 {
			continue
		}

		// ---------------------------------------
		// Create a pending transaction (NO UTXO change yet)
		// ---------------------------------------
		_, _, err = config.Firestore.Collection("pending_transactions").Add(ctx, map[string]interface{}{
			"sender":    walletID,
			"receiver":  ZAKAT_POOL,
			"amount":    zakatAmount,
			"note":      "Monthly zakat deduction",
			"timestamp": time.Now().Unix(),
			"signature": "",              // SYSTEM TRANSACTION
			"inputs":    []interface{}{}, // will be filled during mining
		})

		if err == nil {
			database.AddLog("info", fmt.Sprintf("Zakat pending tx created for wallet %s, amount %.4f", walletID, zakatAmount))
		}
	}

	database.AddLog("success", "Zakat cycle completed — pending transactions created")

	return nil
}
