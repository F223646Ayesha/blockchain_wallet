package database

import "crypto-wallet/config"

func AddPendingTx(data map[string]interface{}) error {
	_, _, err := config.Firestore.Collection("pending_tx").Add(nil, data)
	return err
}
