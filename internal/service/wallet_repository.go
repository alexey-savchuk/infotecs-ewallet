package service

import (
	"context"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
)

//go:generate mockgen -source=wallet_repository.go -destination=mock/wallet_repository_mock.go -package service_mock

type WalletRepository interface {
	Create(ctx context.Context) (*repository.DBWallet, error)
	GetByWalletID(ctx context.Context, walletID string) (*repository.DBWallet, error)
}
