package repository

import (
	"time"

	"github.com/shopspring/decimal"
)

type DBTransfer struct {
	TransferID string
	Time       time.Time
	FromWallet string
	ToWallet   string
	Amount     decimal.Decimal
}
