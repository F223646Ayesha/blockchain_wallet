package blockchain

import (
	"crypto-wallet/utils"
	"fmt"
	"strings"
)

const Difficulty = 5

func MineBlock(b *Block) {
	target := strings.Repeat("0", Difficulty)

	for {
		data := fmt.Sprintf("%d%d%s%s%d",
			b.Index,
			b.Timestamp,
			b.PreviousHash,
			b.MerkleRoot,
			b.Nonce,
		)

		hash := utils.SHA256(data)

		if strings.HasPrefix(hash, target) {
			b.Hash = hash
			return
		}

		b.Nonce++
	}
}
