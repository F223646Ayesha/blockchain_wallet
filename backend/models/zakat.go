package models

type Zakat struct {
	WalletID string  `json:"wallet_id"`
	Amount   float64 `json:"amount"`
	Month    string  `json:"month"`
}
