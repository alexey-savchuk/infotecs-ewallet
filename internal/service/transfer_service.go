package service

import (
	"context"
	"time"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/customerrors"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/handlers"
	"github.com/shopspring/decimal"
)

type TransferService struct {
	repo TransferRepository
}

func NewTransferService(repo TransferRepository) *TransferService {
	return &TransferService{repo: repo}
}

func (ts *TransferService) Create(ctx context.Context, transferDTO *handlers.TransferDTO) (*handlers.TransferDTO, error) {
	if transferDTO.FromWallet == transferDTO.ToWallet {
		return nil, customerrors.ErrSelfTransfer
	}
	if transferDTO.Amount.LessThanOrEqual(decimal.Zero) {
		return nil, customerrors.ErrInvalidAmount
	}

	dbTransfer, err := ts.repo.Create(ctx, &Transfer{
		FromWallet: transferDTO.FromWallet,
		ToWallet:   transferDTO.ToWallet,
		Amount:     transferDTO.Amount,
	})
	if err != nil {
		return nil, err
	}

	return &handlers.TransferDTO{
		Time:       dbTransfer.Time.Format(time.RFC3339),
		FromWallet: dbTransfer.FromWallet,
		ToWallet:   dbTransfer.ToWallet,
		Amount:     dbTransfer.Amount,
	}, nil
}

func (ts *TransferService) GetWalletTransfers(ctx context.Context, walletID string) ([]handlers.TransferDTO, error) {
	dbTransfers, err := ts.repo.GetAllByWalletID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	transferDTOs := make([]handlers.TransferDTO, 0)
	for _, dbTransfer := range dbTransfers {
		transferDTOs = append(transferDTOs, handlers.TransferDTO{
			Time:       dbTransfer.Time.Format(time.RFC3339),
			FromWallet: dbTransfer.FromWallet,
			ToWallet:   dbTransfer.ToWallet,
			Amount:     dbTransfer.Amount,
		})
	}
	return transferDTOs, nil
}
