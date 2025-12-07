package blockchain

import (
	"crypto-wallet/utils"
	"fmt"
)

// RecomputeHash rebuilds the block hash using the SAME formula as mining.
func RecomputeHash(b *Block) string {
	data := fmt.Sprintf("%d%d%s%s%d",
		b.Index,
		b.Timestamp,
		b.PreviousHash,
		b.MerkleRoot,
		b.Nonce,
	)

	return utils.SHA256(data)
}
