package service

import (
	"github.com/shopspring/decimal"
)

type Transfer struct {
	FromWallet string
	ToWallet   string
	Amount     decimal.Decimal
}
