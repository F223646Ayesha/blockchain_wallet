package blockchain

type Block struct {
	Index        int                      `json:"index" firestore:"index"`
	Timestamp    int64                    `json:"timestamp" firestore:"timestamp"`
	Transactions []map[string]interface{} `json:"transactions" firestore:"transactions"`
	PreviousHash string                   `json:"previous_hash" firestore:"previous_hash"`
	Nonce        int                      `json:"nonce" firestore:"nonce"`
	Hash         string                   `json:"hash" firestore:"hash"`
	MerkleRoot   string                   `json:"merkle_root" firestore:"merkle_root"`
}
