package database

import "crypto-wallet/config"

func SaveBlock(data map[string]interface{}) error {
	_, _, err := config.Firestore.Collection("blocks").Add(nil, data)
	return err
}
