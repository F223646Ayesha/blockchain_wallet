package models

type Wallet struct {
	WalletID      string   `json:"wallet_id" firestore:"wallet_id"`
	PublicKey     string   `json:"public_key" firestore:"public_key"`
	PrivateKeyEnc string   `json:"private_key_enc" firestore:"private_key_enc"`
	Beneficiaries []string `json:"beneficiaries" firestore:"beneficiaries"`
}
