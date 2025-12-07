package services

import (
	"context"
	"crypto-wallet/config"
)

type SystemAnalytics struct {
	TotalWallets        int     `json:"total_wallets"`
	BlocksMined         int     `json:"blocks_mined"`
	PendingTransactions int     `json:"pending_transactions"`
	CompletedTx         int     `json:"completed_transactions"`
	TotalZakat          float64 `json:"total_zakat_collected"`
}

func GetSystemAnalyticsService() (*SystemAnalytics, error) {
	ctx := context.Background()

	result := &SystemAnalytics{}

	// 1. Count wallets
	wallets, err := config.Firestore.Collection("wallets").Documents(ctx).GetAll()
	if err == nil {
		result.TotalWallets = len(wallets)
	}

	// 2. Count blocks (mined)
	blocks, err := config.Firestore.Collection("blocks").Documents(ctx).GetAll()
	if err == nil {
		result.BlocksMined = len(blocks)
	}

	// 3. Count pending transactions
	pend, err := config.Firestore.Collection("pending_transactions").Documents(ctx).GetAll()
	if err == nil {
		result.PendingTransactions = len(pend)
	}

	// 4. Count completed transactions
	txs, err := config.Firestore.Collection("transactions").Documents(ctx).GetAll()
	if err == nil {
		result.CompletedTx = len(txs)
	}

	// 5. Sum zakat collected
	zakatTxs, err := config.Firestore.Collection("transactions").
		Where("type", "==", "zakat_deduction").Documents(ctx).GetAll()
	if err == nil {
		sum := 0.0
		for _, d := range zakatTxs {
			if amt, ok := d.Data()["amount"].(float64); ok {
				sum += amt
			}
		}
		result.TotalZakat = sum
	}

	return result, nil
}
