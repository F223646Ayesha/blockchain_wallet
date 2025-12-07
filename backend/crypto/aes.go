package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// ------------------------------------------
// Derive key & IV exactly like OpenSSL
// ------------------------------------------
func deriveKey(password, salt []byte) (key, iv []byte) {
	prev := []byte{}
	out := []byte{}

	for len(out) < 48 { // 32-byte key + 16-byte iv
		h := sha256.New()
		h.Write(prev)
		h.Write(password)
		h.Write(salt)
		prev = h.Sum(nil)
		out = append(out, prev...)
	}

	key = out[:32]
	iv = out[32:48]
	return
}

// ------------------------------------------
// AES-256-CBC Encryption (OpenSSL compatible)
// ------------------------------------------
func EncryptAES(password []byte, plaintext string) (string, error) {
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key, iv := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	plain := pkcs7Pad([]byte(plaintext), aes.BlockSize)
	ciphertext := make([]byte, len(plain))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plain)

	final := append([]byte("Salted__"), salt...)
	final = append(final, ciphertext...)

	return base64.StdEncoding.EncodeToString(final), nil
}

// ------------------------------------------
// AES-256-CBC Decryption
// ------------------------------------------
func DecryptAES(password []byte, b64 string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	if len(raw) < 16 || string(raw[:8]) != "Salted__" {
		return "", errors.New("invalid data")
	}

	salt := raw[8:16]
	ciphertext := raw[16:]

	key, iv := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext not multiple of block size")
	}

	plain := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plain, ciphertext)

	plain, err = pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(plain), nil
}

// ------------------------------------------
// PKCS7 padding
// ------------------------------------------
func pkcs7Pad(b []byte, blockSize int) []byte {
	pad := blockSize - (len(b) % blockSize)
	padding := bytes.Repeat([]byte{byte(pad)}, pad)
	return append(b, padding...)
}

func pkcs7Unpad(b []byte, blockSize int) ([]byte, error) {
	l := len(b)
	if l == 0 || l%blockSize != 0 {
		return nil, errors.New("invalid padding")
	}
	pad := int(b[l-1])
	if pad == 0 || pad > blockSize {
		return nil, errors.New("invalid padding size")
	}
	for i := 0; i < pad; i++ {
		if b[l-1-i] != byte(pad) {
			return nil, errors.New("bad padding")
		}
	}
	return b[:l-pad], nil
}
