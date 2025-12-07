package database

import (
	"context"
	"crypto-wallet/config"
)

func GetFullTransactionHistory(walletID string) ([]map[string]interface{}, error) {
	ctx := context.Background()

	// 1️⃣ Read all blocks from Firestore
	iter := config.Firestore.Collection("blocks").Documents(ctx)

	var fullHistory []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		block := doc.Data()
		txs := block["transactions"].([]interface{})

		// 2️⃣ Loop through each tx inside this block
		for _, t := range txs {
			tx := t.(map[string]interface{})

			sender := tx["sender"]
			receiver := tx["receiver"]

			// Only include transactions related to this wallet
			if sender == walletID || receiver == walletID {
				fullHistory = append(fullHistory, tx)
			}
		}
	}

	return fullHistory, nil
}
