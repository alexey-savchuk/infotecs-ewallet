package service

import (
	"context"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
)

//go:generate mockgen -source=transfer_repository.go -destination=mock/transfer_repository_mock.go -package service_mock

type TransferRepository interface {
	Create(ctx context.Context, transfer *Transfer) (*repository.DBTransfer, error)
	GetAllByWalletID(ctx context.Context, walletID string) ([]*repository.DBTransfer, error)
}
