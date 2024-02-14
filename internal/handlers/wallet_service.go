package handlers

import "context"

type WalletService interface {
	Create(ctx context.Context) (*WalletDTO, error)
	Get(ctx context.Context, walletID string) (*WalletDTO, error)
}
