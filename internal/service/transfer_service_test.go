package service_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/customerrors"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/handlers"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/service"
	service_mock "github.com/alexey-savchuk/infotecs-ewallet/internal/service/mock"
	"github.com/shopspring/decimal"
	"go.uber.org/mock/gomock"
)

func TestTransferService_Create(t *testing.T) {
	t.Run("SuccessfulTransfer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		dbTransfer := repository.DBTransfer{
			TransferID: "1",
			FromWallet: "1",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}
		transfer := service.Transfer{
			FromWallet: dbTransfer.FromWallet,
			ToWallet:   dbTransfer.ToWallet,
			Amount:     dbTransfer.Amount,
		}
		transferDTO := handlers.TransferDTO{
			FromWallet: transfer.FromWallet,
			ToWallet:   transfer.ToWallet,
			Amount:     transfer.Amount,
		}

		repo := service_mock.NewMockTransferRepository(ctrl)
		repo.
			EXPECT().
			Create(context.Background(), &transfer).
			Return(&dbTransfer, nil)

		ws := service.NewTransferService(repo)
		result, err := ws.Create(context.Background(), &transferDTO)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if result == nil || *result != transferDTO {
			t.Errorf("got %v, want %v", result, transferDTO)
		}
	})

	t.Run("SelfTransfer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		transferDTO := handlers.TransferDTO{
			FromWallet: "1",
			ToWallet:   "1",
			Amount:     decimal.NewFromFloat(100.0),
		}

		repo := service_mock.NewMockTransferRepository(ctrl)

		ws := service.NewTransferService(repo)
		result, err := ws.Create(context.Background(), &transferDTO)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrSelfTransfer {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrSelfTransfer)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})

	t.Run("InvalidAmount_Negative", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		transferDTO := handlers.TransferDTO{
			FromWallet: "1",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(-100.0),
		}

		repo := service_mock.NewMockTransferRepository(ctrl)

		ws := service.NewTransferService(repo)
		result, err := ws.Create(context.Background(), &transferDTO)
		if err == nil {
			t.Error("expected error")
		}
		if err != customerrors.ErrInvalidAmount {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrInvalidAmount)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})

	t.Run("InvalidAmount_Zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		transferDTO := handlers.TransferDTO{
			FromWallet: "1",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(0.0),
		}

		repo := service_mock.NewMockTransferRepository(ctrl)

		ws := service.NewTransferService(repo)
		result, err := ws.Create(context.Background(), &transferDTO)
		if err == nil {
			t.Error("expected error")
		}
		if err != customerrors.ErrInvalidAmount {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrInvalidAmount)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})

	t.Run("FromWalletNotExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		transfer := service.Transfer{
			FromWallet: "NotExistingID",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}
		transferDTO := handlers.TransferDTO{
			FromWallet: transfer.FromWallet,
			ToWallet:   transfer.ToWallet,
			Amount:     transfer.Amount,
		}

		repo := service_mock.NewMockTransferRepository(ctrl)
		repo.
			EXPECT().
			Create(context.Background(), &transfer).
			Return(nil, customerrors.ErrFromWalletNotExists)

		ws := service.NewTransferService(repo)
		result, err := ws.Create(context.Background(), &transferDTO)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrFromWalletNotExists {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrFromWalletNotExists)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})

	t.Run("ToWalletNotExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		transfer := service.Transfer{
			FromWallet: "1",
			ToWallet:   "NotExistingID",
			Amount:     decimal.NewFromFloat(100.0),
		}
		transferDTO := handlers.TransferDTO{
			FromWallet: transfer.FromWallet,
			ToWallet:   transfer.ToWallet,
			Amount:     transfer.Amount,
		}

		repo := service_mock.NewMockTransferRepository(ctrl)
		repo.
			EXPECT().
			Create(context.Background(), &transfer).
			Return(nil, customerrors.ErrToWalletNotExists)

		ws := service.NewTransferService(repo)
		result, err := ws.Create(context.Background(), &transferDTO)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrToWalletNotExists {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrToWalletNotExists)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})

	t.Run("InsufficientFunds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		transfer := service.Transfer{
			FromWallet: "1",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}
		transferDTO := handlers.TransferDTO{
			FromWallet: transfer.FromWallet,
			ToWallet:   transfer.ToWallet,
			Amount:     transfer.Amount,
		}

		repo := service_mock.NewMockTransferRepository(ctrl)
		repo.
			EXPECT().
			Create(context.Background(), &transfer).
			Return(nil, customerrors.ErrInsufficientFunds)

		ws := service.NewTransferService(repo)
		result, err := ws.Create(context.Background(), &transferDTO)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrInsufficientFunds {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrInsufficientFunds)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})
}

func TestTransferService_GetWalletTransfers(t *testing.T) {
	t.Run("WalletExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		walletID := "1"
		dbTransfers := []repository.DBTransfer{
			{
				TransferID: "1",
				Time:       time.Now().Add(time.Hour),
				FromWallet: walletID,
				ToWallet:   "2",
				Amount:     decimal.NewFromFloat(100.0),
			},
			{
				TransferID: "2",
				Time:       time.Now(),
				FromWallet: "3",
				ToWallet:   walletID,
				Amount:     decimal.NewFromFloat(200.0),
			},
		}
		transferDTOs := []handlers.TransferDTO{
			{
				Time:       dbTransfers[0].Time.Format(time.RFC3339),
				FromWallet: dbTransfers[0].FromWallet,
				ToWallet:   dbTransfers[0].ToWallet,
				Amount:     dbTransfers[0].Amount,
			},
			{
				Time:       dbTransfers[1].Time.Format(time.RFC3339),
				FromWallet: dbTransfers[1].FromWallet,
				ToWallet:   dbTransfers[1].ToWallet,
				Amount:     dbTransfers[1].Amount,
			},
		}

		repo := service_mock.NewMockTransferRepository(ctrl)
		repo.
			EXPECT().
			GetAllByWalletID(context.Background(), walletID).
			Return(dbTransfers, nil)

		ws := service.NewTransferService(repo)
		result, err := ws.GetWalletTransfers(context.Background(), walletID)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if result == nil || !reflect.DeepEqual(result, transferDTOs) {
			t.Errorf("got %v, want %v", result, transferDTOs)
		}
	})

	t.Run("WalletNotExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		walletID := "1"

		repo := service_mock.NewMockTransferRepository(ctrl)
		repo.
			EXPECT().
			GetAllByWalletID(context.Background(), walletID).
			Return(nil, customerrors.ErrWalletNotExists)

		ws := service.NewTransferService(repo)
		result, err := ws.GetWalletTransfers(context.Background(), walletID)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrWalletNotExists {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrWalletNotExists)
		}
		if result != nil {
			t.Errorf("got %v, want %v", result, nil)
		}
	})
}
