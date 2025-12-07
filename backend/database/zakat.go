package database

import "crypto-wallet/config"

func SaveZakat(data map[string]interface{}) error {
	_, _, err := config.Firestore.Collection("zakat").Add(nil, data)
	return err
}
