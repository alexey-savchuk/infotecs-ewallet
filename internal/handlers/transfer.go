package handlers

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransferDTO struct {
	Time       time.Time       `json:"time"`
	FromWallet string          `json:"from"`
	ToWallet   string          `json:"to"`
	Amount     decimal.Decimal `json:"amount"`
}
