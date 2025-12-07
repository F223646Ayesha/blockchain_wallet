package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"math/big"
)

type ECDSASignature struct {
	R, S *big.Int
}

// SignPayload signs ANY string using the private key D value (hex string)
func SignPayload(privateKey *ecdsa.PrivateKey, payload string) ([]byte, error) {
	hash := sha256.Sum256([]byte(payload))

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(ECDSASignature{r, s})
}

// VerifySignature verifies the signature using public key bytes (X+Y hex)
func VerifySignature(publicKey ecdsa.PublicKey, payload string, signature []byte) bool {
	var sig ECDSASignature

	_, err := asn1.Unmarshal(signature, &sig)
	if err != nil {
		return false
	}

	hash := sha256.Sum256([]byte(payload))

	return ecdsa.Verify(&publicKey, hash[:], sig.R, sig.S)
}
