package service

import (
	"context"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
)

type TransferRepository interface {
	Create(ctx context.Context, transfer *Transfer) (*repository.DBTransfer, error)
	GetAllByWalletID(ctx context.Context, walletID string) ([]*repository.DBTransfer, error)
}
