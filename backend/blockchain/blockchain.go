package blockchain

func NewBlock(index int, timestamp int64, prevHash string, txs []map[string]interface{}) *Block {

	// Mining reward transaction placeholder
	rewardTx := map[string]interface{}{
		"sender":    "SYSTEM",
		"receiver":  "", // to be filled at mining time
		"amount":    50,
		"note":      "mining_reward",
		"timestamp": timestamp,
	}

	// reward always at index 0
	updatedTx := append([]map[string]interface{}{rewardTx}, txs...)

	b := &Block{
		Index:        index,
		Timestamp:    timestamp,
		PreviousHash: prevHash,
		Transactions: updatedTx,
		Nonce:        0,
	}

	b.MerkleRoot = ComputeMerkleRoot(updatedTx)

	return b
}

// Called before mining to assign miner
func (b *Block) AssignMiner(wallet string) {
	b.Transactions[0]["receiver"] = wallet
}
