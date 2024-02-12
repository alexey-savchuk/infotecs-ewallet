package model

import "github.com/shopspring/decimal"

type Wallet struct {
	WalletID string
	Balance  decimal.Decimal
}
