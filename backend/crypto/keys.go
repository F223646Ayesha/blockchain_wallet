package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"math/big"
)

type KeyPair struct {
	PrivateKey      *ecdsa.PrivateKey
	PrivateKeyPKCS8 []byte
	PrivateKeyHex   string
	PublicKey       []byte
	PublicKeyHex    string
}

func GenerateKeyPair() (*KeyPair, error) {

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	pub := append([]byte{0x04},
		append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)...,
	)

	pkcs8, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		PrivateKey:      priv,
		PrivateKeyPKCS8: pkcs8,
		PrivateKeyHex:   hex.EncodeToString(priv.D.Bytes()),
		PublicKey:       pub,
		PublicKeyHex:    hex.EncodeToString(pub),
	}, nil
}

func GenerateWalletID(pub []byte) string {
	hash := sha256.Sum256(pub)
	return hex.EncodeToString(hash[:])
}

// ======================================================
// PRIVATE KEY FROM HEX (D → full ECDSA private key)
// ======================================================

func PrivateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {

	dBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}

	if len(dBytes) != 32 {
		return nil, errors.New("invalid private key length: must be 32 bytes")
	}

	curve := elliptic.P256()

	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(dBytes)

	// Recompute public key from private key
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(dBytes)

	return priv, nil
}

// ======================================================
// PUBLIC KEY FROM HEX (SEC1: 04 | X | Y)
// ======================================================

func PublicKeyFromHex(pubHex string) (*ecdsa.PublicKey, error) {

	pubBytes, err := hex.DecodeString(pubHex)
	if err != nil {
		return nil, err
	}

	if len(pubBytes) != 65 {
		return nil, errors.New("invalid public key length: must be 65 bytes")
	}

	if pubBytes[0] != 0x04 {
		return nil, errors.New("invalid SEC1 format: expected 0x04 prefix")
	}

	curve := elliptic.P256()

	x := new(big.Int).SetBytes(pubBytes[1:33])
	y := new(big.Int).SetBytes(pubBytes[33:65])

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}
