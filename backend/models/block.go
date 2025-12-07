package models

type Block struct {
	Index        int         `json:"index"`
	Timestamp    int64       `json:"timestamp"`
	PreviousHash string      `json:"previous_hash"`
	Hash         string      `json:"hash"`
	Nonce        int         `json:"nonce"`
	Transactions interface{} `json:"transactions"`
}
