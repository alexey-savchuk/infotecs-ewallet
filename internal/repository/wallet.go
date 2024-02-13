package repository

import "github.com/shopspring/decimal"

type DBWallet struct {
	WalletID string
	Balance  decimal.Decimal
}
