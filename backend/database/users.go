package database

import "crypto-wallet/config"

func CreateUser(data map[string]interface{}) error {
	_, _, err := config.Firestore.Collection("users").Add(nil, data)
	return err
}
