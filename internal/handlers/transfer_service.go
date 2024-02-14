package handlers

import "context"

type TransferService interface {
	Create(ctx context.Context, transfer *TransferDTO) (*TransferDTO, error)
	GetWalletTransfers(ctx context.Context, walletID string) ([]*TransferDTO, error)
}
