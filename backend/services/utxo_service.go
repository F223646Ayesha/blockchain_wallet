package services

import (
	"crypto-wallet/database"
	"crypto-wallet/models"
	"crypto-wallet/utils"
)

// GetBalance calculates wallet balance using UTXOs
func GetBalance(walletID string) (float64, error) {

	utxos, err := database.GetUnspentUTXOs(walletID)
	if err != nil {
		return 0, err
	}

	total := 0.0
	for _, u := range utxos {
		total += u.Amount
	}

	return total, nil
}

// SelectUTXOs picks enough UTXOs to cover a transaction amount
func SelectUTXOs(walletID string, amount float64) ([]models.UTXO, float64, error) {

	// Get all usable UTXOs
	utxos, err := database.GetUnspentUTXOs(walletID)
	if err != nil {
		return nil, 0, err
	}

	var selected []models.UTXO
	sum := 0.0

	// Bitcoin-style greedy selection
	for _, u := range utxos {
		selected = append(selected, u)
		sum += u.Amount

		if sum >= amount {
			change := sum - amount
			return selected, change, nil
		}
	}

	// Not enough funds
	return nil, 0, utils.ErrInsufficientBalance
}
