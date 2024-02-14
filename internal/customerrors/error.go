package customerrors

import "errors"

var (
	ErrSelfTransfer        = errors.New("self transfer")
	ErrWalletNotExists     = errors.New("wallet not exists")
	ErrFromWalletNotExists = errors.New("from wallet not exists")
	ErrToWalletNotExists   = errors.New("to wallet not exists")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrInvalidAmount       = errors.New("invalid amount")
)
