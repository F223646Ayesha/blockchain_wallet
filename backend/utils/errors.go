package utils

import "errors"

var (
	ErrInvalidWallet       = errors.New("Invalid wallet ID")
	ErrInvalidSignature    = errors.New("Digital signature verification failed")
	ErrInsufficientBalance = errors.New("Insufficient balance")
	ErrDoubleSpend         = errors.New("Double spend detected")
)
