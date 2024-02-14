package handlers

import "github.com/shopspring/decimal"

type WalletDTO struct {
	WalletID string          `json:"id" param:"walletId"`
	Balance  decimal.Decimal `json:"balance"`
}
