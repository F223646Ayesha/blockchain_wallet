package services

import (
	"crypto-wallet/blockchain"
	"fmt"
)

func ValidateBlockchainService(blocks []blockchain.Block) (bool, string) {

	// Validate genesis block first
	if len(blocks) == 0 {
		return false, "No blocks found"
	}

	for i := 0; i < len(blocks); i++ {
		current := blocks[i]

		// ------------------------
		// Rule 1: Recalculate Hash
		// ------------------------
		expectedHash := blockchain.RecomputeHash(&current)

		if current.Hash != expectedHash {
			return false, fmt.Sprintf("❌ Hash mismatch at block %d", current.Index)
		}

		// ------------------------
		// Rule 2: Validate Previous Hash Link
		// ------------------------
		if i > 0 && current.PreviousHash != blocks[i-1].Hash {
			return false, fmt.Sprintf("❌ Previous hash mismatch at block %d", current.Index)
		}

		// ------------------------
		// Rule 3: Validate Proof-of-Work
		// ------------------------
		diff := blockchain.Difficulty
		prefix := expectedHash[:diff]

		for j := 0; j < diff; j++ {
			if prefix[j] != '0' {
				return false, fmt.Sprintf("❌ Invalid PoW at block %d", current.Index)
			}
		}
	}

	return true, "✅ Blockchain is valid"
}
