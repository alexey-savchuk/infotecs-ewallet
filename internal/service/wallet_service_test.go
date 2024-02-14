package service_test

import (
	"context"
	"testing"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/customerrors"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/handlers"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/service"
	service_mock "github.com/alexey-savchuk/infotecs-ewallet/internal/service/mock"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"
)

func TestWalletService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbWallet := repository.DBWallet{
		WalletID: "1",
		Balance:  decimal.NewFromFloat(100.0),
	}
	walletDTO := handlers.WalletDTO{
		WalletID: dbWallet.WalletID,
		Balance:  dbWallet.Balance,
	}

	repo := service_mock.NewMockWalletRepository(ctrl)
	repo.
		EXPECT().
		Create(context.Background()).
		Return(&dbWallet, nil)

	ws := service.NewWalletService(repo)
	wallet, err := ws.Create(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if wallet == nil || *wallet != walletDTO {
		t.Errorf("got %v, want %v", wallet, walletDTO)
	}
}

func TestWalletService_Get(t *testing.T) {
	t.Run("WalletExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbWallet := repository.DBWallet{
			WalletID: "1",
			Balance:  decimal.NewFromFloat(100.0),
		}
		walletDTO := handlers.WalletDTO{
			WalletID: dbWallet.WalletID,
			Balance:  dbWallet.Balance,
		}

		repo := service_mock.NewMockWalletRepository(ctrl)
		repo.
			EXPECT().
			GetByWalletID(context.Background(), dbWallet.WalletID).
			Return(&dbWallet, nil)

		ws := service.NewWalletService(repo)
		wallet, err := ws.Get(context.Background(), dbWallet.WalletID)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if wallet == nil || *wallet != walletDTO {
			t.Errorf("got %v, want %v", wallet, walletDTO)
		}
	})

	t.Run("WalletNotExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		walletID := "NotExistingID"
		repo := service_mock.NewMockWalletRepository(ctrl)
		repo.
			EXPECT().
			GetByWalletID(context.Background(), walletID).
			Return(nil, customerrors.ErrWalletNotExists)

		ws := service.NewWalletService(repo)
		wallet, err := ws.Get(context.Background(), walletID)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrWalletNotExists {
			t.Errorf("got %v, want %v", err, customerrors.ErrWalletNotExists)
		}
		if wallet != nil {
			t.Errorf("got %v, want %v", wallet, nil)
		}
	})
}
