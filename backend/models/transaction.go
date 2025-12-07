package models

type Transaction struct {
	WalletID     string  `json:"wallet_id"`
	Type         string  `json:"type"` // sent / received / zakat
	Counterparty string  `json:"counterparty"`
	Amount       float64 `json:"amount"`
	Note         string  `json:"note"`
	Timestamp    int64   `json:"timestamp"`
}
