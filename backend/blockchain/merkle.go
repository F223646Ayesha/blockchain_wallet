package blockchain

import (
	"crypto-wallet/utils"
	"encoding/json"
)

// ComputeMerkleRoot builds a simple Merkle root from transaction data.
func ComputeMerkleRoot(txs []map[string]interface{}) string {
	if len(txs) == 0 {
		return ""
	}

	// Step 1: hash each transaction as a leaf
	leaves := []string{}
	for _, tx := range txs {
		b, _ := json.Marshal(tx)
		leaves = append(leaves, utils.SHA256(string(b)))
	}

	// If only 1 transaction, its hash is the merkle root
	if len(leaves) == 1 {
		return leaves[0]
	}

	// Step 2: repeatedly hash pairs until one root remains
	for len(leaves) > 1 {
		var next []string

		for i := 0; i < len(leaves); i += 2 {
			if i+1 < len(leaves) {
				next = append(next, utils.SHA256(leaves[i]+leaves[i+1]))
			} else {
				// if odd count, duplicate last node
				next = append(next, utils.SHA256(leaves[i]+leaves[i]))
			}
		}

		leaves = next
	}

	return leaves[0]
}
