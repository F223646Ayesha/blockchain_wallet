package database

import (
	"context" // ⭐ REQUIRED
	"crypto-wallet/config"
	"crypto-wallet/models"

	"cloud.google.com/go/firestore"
)

var UTXO_COLLECTION = "utxos"

// --------------------------------------------------
// Create UTXO
// --------------------------------------------------
func AddUTXO(utxo models.UTXO) error {
	ctx := context.Background() // ⭐ FIX
	_, _, err := config.Firestore.Collection(UTXO_COLLECTION).Add(ctx, utxo)
	return err
}

// --------------------------------------------------
// Get all unspent UTXOs for a wallet
// --------------------------------------------------
func GetUnspentUTXOs(walletID string) ([]models.UTXO, error) {

	ctx := context.Background() // ⭐ FIX

	q := config.Firestore.Collection(UTXO_COLLECTION).
		Where("wallet_id", "==", walletID).
		Where("spent", "==", false)

	docs, err := q.Documents(ctx).GetAll() // ⭐ FIX
	if err != nil {
		return nil, err
	}

	utxos := []models.UTXO{}
	for _, doc := range docs {
		var u models.UTXO
		doc.DataTo(&u)
		u.ID = doc.Ref.ID
		utxos = append(utxos, u)
	}

	return utxos, nil
}

// --------------------------------------------------
// Mark UTXO as spent
// --------------------------------------------------
func MarkUTXOSpent(id string) error {
	ctx := context.Background() // ⭐ FIX

	_, err := config.Firestore.Collection(UTXO_COLLECTION).
		Doc(id).
		Update(ctx, []firestore.Update{
			{Path: "spent", Value: true},
		})
	return err
}
