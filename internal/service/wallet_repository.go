package service

import (
	"context"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
)

type WalletRepository interface {
	Create(ctx context.Context) (*repository.DBWallet, error)
	GetByWalletID(ctx context.Context, walletID string) (*repository.DBWallet, error)
}
