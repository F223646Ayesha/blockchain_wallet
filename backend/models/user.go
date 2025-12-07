package models

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	CNIC     string `json:"cnic"`
	WalletID string `json:"wallet_id"`
}
