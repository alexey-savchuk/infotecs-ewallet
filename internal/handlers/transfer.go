package handlers

import (
	"github.com/shopspring/decimal"
)

type TransferDTO struct {
	Time       string          `json:"time"`
	FromWallet string          `json:"from" param:"walletId"`
	ToWallet   string          `json:"to"`
	Amount     decimal.Decimal `json:"amount"`
}
