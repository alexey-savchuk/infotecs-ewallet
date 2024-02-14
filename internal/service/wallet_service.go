package service

import (
	"context"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/handlers"
)

type WalletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (ws *WalletService) Create(ctx context.Context) (*handlers.WalletDTO, error) {
	dbWallet, err := ws.repo.Create(ctx)
	if err != nil {
		return nil, err
	}
	return &handlers.WalletDTO{
		WalletID: dbWallet.WalletID,
		Balance:  dbWallet.Balance,
	}, nil
}

func (ws *WalletService) Get(ctx context.Context, walletID string) (*handlers.WalletDTO, error) {
	dbWallet, err := ws.repo.GetByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}
	return &handlers.WalletDTO{
		WalletID: dbWallet.WalletID,
		Balance:  dbWallet.Balance,
	}, nil
}
