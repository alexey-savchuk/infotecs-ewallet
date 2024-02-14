package handlers

import "github.com/shopspring/decimal"

type WalletDTO struct {
	WalletID string          `json:"id" query:"walletId"`
	Balance  decimal.Decimal `json:"balance"`
}
