package services

import (
	"context"
	"crypto-wallet/config"
	"errors"
	"time"
)

type MonthlySummary struct {
	Month    string  `json:"month"`
	Sent     float64 `json:"sent"`
	Received float64 `json:"received"`
	Zakat    float64 `json:"zakat"`
	TxCount  int     `json:"txCount"`
}

type WalletReport struct {
	TotalSent        float64          `json:"totalSent"`
	TotalReceived    float64          `json:"totalReceived"`
	ZakatPaid        float64          `json:"zakatPaid"`
	SentCount        int              `json:"sentCount"`
	ReceivedCount    int              `json:"receivedCount"`
	TransactionCount int              `json:"transactionCount"`
	BlocksMined      int              `json:"blocksMined"`
	MonthlySummary   []MonthlySummary `json:"monthlySummary"`
}

// ======================================================
// GENERATE WALLET REPORT
// ======================================================
func GenerateWalletReport(walletId string) (*WalletReport, error) {

	if walletId == "" {
		return nil, errors.New("wallet id required")
	}

	ctx := context.Background()

	// ============================
	// 1. FETCH ALL TRANSACTIONS
	// ============================
	txDocs, err := config.Firestore.Collection("transactions").
		Where("wallet_id", "==", walletId).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	report := &WalletReport{
		MonthlySummary: []MonthlySummary{},
	}

	// month → summary
	monthlyMap := make(map[string]*MonthlySummary)

	// Loop through transactions
	for _, doc := range txDocs {
		data := doc.Data()

		amount := data["amount"].(float64)
		txType := data["type"].(string)
		timestamp := data["timestamp"].(int64)

		monthStr := time.Unix(timestamp, 0).Format("2006-01")

		// ensure month exists
		if _, ok := monthlyMap[monthStr]; !ok {
			monthlyMap[monthStr] = &MonthlySummary{
				Month:    monthStr,
				Sent:     0,
				Received: 0,
				Zakat:    0,
				TxCount:  0,
			}
		}

		// update counts
		monthlyMap[monthStr].TxCount++

		// update summary numbers
		switch txType {
		case "sent":
			report.TotalSent += amount
			report.SentCount += 1
			monthlyMap[monthStr].Sent += amount

		case "received":
			report.TotalReceived += amount
			report.ReceivedCount += 1
			monthlyMap[monthStr].Received += amount

		case "zakat_deduction":
			report.ZakatPaid += amount
			monthlyMap[monthStr].Zakat += amount
		}
	}

	report.TransactionCount = len(txDocs)

	// Convert monthlyMap → list
	for _, m := range monthlyMap {
		report.MonthlySummary = append(report.MonthlySummary, *m)
	}

	// ============================
	// 2. BLOCKS MINED BY USER
	// ============================
	blockDocs, err := config.Firestore.Collection("blocks").
		Where("miner", "==", walletId).
		Documents(ctx).
		GetAll()

	if err == nil {
		report.BlocksMined = len(blockDocs)
	}

	return report, nil
}
