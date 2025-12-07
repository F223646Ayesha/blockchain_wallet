package models

type UTXO struct {
	ID       string  `json:"id" firestore:"id"`
	WalletID string  `json:"wallet_id" firestore:"wallet_id"`
	Amount   float64 `json:"amount" firestore:"amount"`
	Spent    bool    `json:"spent" firestore:"spent"`
	TxID     string  `json:"tx_id" firestore:"tx_id"`
}
