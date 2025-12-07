package database

import "crypto-wallet/config"

func CreateWallet(data map[string]interface{}) error {
	_, _, err := config.Firestore.Collection("wallets").Add(nil, data)
	return err
}
