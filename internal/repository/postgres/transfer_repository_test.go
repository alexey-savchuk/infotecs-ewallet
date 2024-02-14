package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/customerrors"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository/postgres"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/service"
	"github.com/shopspring/decimal"
)

func TestTransferRepository_Create(t *testing.T) {
	t.Run("WalletsExist_EnoughFunds", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		transfer := &service.Transfer{
			FromWallet: "1",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}

		mock.ExpectBegin()
		query := `SELECT balance FROM wallet WHERE wallet_id = \$1`
		mock.
			ExpectQuery(query).
			WithArgs(transfer.FromWallet).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200.0))
		query = `UPDATE wallet SET balance = balance \- \$1 WHERE wallet_id = \$2`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.FromWallet).
			WillReturnResult(sqlmock.NewResult(1, 1))
		query = `UPDATE wallet SET balance = balance \+ \$1 WHERE wallet_id = \$2`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.ToWallet).
			WillReturnResult(sqlmock.NewResult(1, 1))
		query = `INSERT INTO transfer \(from_wallet, to_wallet, amount\)
				 VALUES \(\$1, \$2, \$3\)
				 RETURNING transfer_id, time, from_wallet, to_wallet, amount`
		mock.
			ExpectQuery(query).
			WithArgs().
			WillReturnRows(
				sqlmock.
					NewRows([]string{"transfer_id", "time", "from_wallet", "to_wallet", "amount"}).
					AddRow("1", time.Now(), transfer.FromWallet, transfer.ToWallet, transfer.Amount),
			)
		mock.ExpectCommit()

		tr := postgres.NewTransferRepository(db)

		_, err = tr.Create(context.Background(), transfer)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("NotEnoughFunds", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		transfer := &service.Transfer{
			FromWallet: "1",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}

		mock.ExpectBegin()
		query := `SELECT balance FROM wallet WHERE wallet_id = \$1`
		mock.
			ExpectQuery(query).
			WithArgs(transfer.FromWallet).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(50.0))
		mock.ExpectRollback()

		tr := postgres.NewTransferRepository(db)

		_, err = tr.Create(context.Background(), transfer)
		if err == nil {
			t.Error("expected error")
		}
		if err != customerrors.ErrInsufficientFunds {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrInsufficientFunds)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("FromWalletNotExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		transfer := &service.Transfer{
			FromWallet: "NotExistingID",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}

		mock.ExpectBegin()
		query := `SELECT balance FROM wallet WHERE wallet_id = \$1`
		mock.
			ExpectQuery(query).
			WithArgs(transfer.FromWallet).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		tr := postgres.NewTransferRepository(db)

		_, err = tr.Create(context.Background(), transfer)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrFromWalletNotExists {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrFromWalletNotExists)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("ToWalletNotExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		transfer := &service.Transfer{
			FromWallet: "1",
			ToWallet:   "NotExistingID",
			Amount:     decimal.NewFromFloat(100.0),
		}

		mock.ExpectBegin()
		query := `SELECT balance FROM wallet WHERE wallet_id = \$1`
		mock.
			ExpectQuery(query).
			WithArgs(transfer.FromWallet).
			WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200.0))
		query = `UPDATE wallet SET balance = balance \- \$1 WHERE wallet_id = \$2`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.FromWallet).
			WillReturnResult(sqlmock.NewResult(1, 1))
		query = `UPDATE wallet SET balance = balance \+ \$1 WHERE wallet_id = \$2`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.ToWallet).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		tr := postgres.NewTransferRepository(db)

		_, err = tr.Create(context.Background(), transfer)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrToWalletNotExists {
			t.Errorf("unexpected error: %s, want %s", err, customerrors.ErrToWalletNotExists)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestTransferRepository_GetAllByWalletID(t *testing.T) {
	t.Run("WalletExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		walletID := "1"

		mock.ExpectBegin()
		query := `SELECT wallet_id FROM wallet WHERE wallet_id = \$1`
		mock.
			ExpectQuery(query).
			WithArgs(walletID).
			WillReturnRows(sqlmock.NewRows([]string{"wallet_id"}).AddRow(walletID))
		query = `SELECT transfer_id, time, from_wallet, to_wallet, amount
				  FROM transfer WHERE from_wallet = \$1 OR to_wallet = \$1 ORDER BY 2 DESC`
		mock.
			ExpectQuery(query).
			WithArgs(walletID).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"transfer_id", "time", "from_wallet", "to_wallet", "amount"}).
					AddRow("1", time.Now(), walletID, "2", 100.0).
					AddRow("2", time.Now(), "3", walletID, 200.0),
			)
		mock.ExpectCommit()

		tr := postgres.NewTransferRepository(db)

		_, err = tr.GetAllByWalletID(context.Background(), walletID)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("WalletNotExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		walletID := "NotExistingID"
		mock.ExpectBegin()
		query := `SELECT wallet_id FROM wallet WHERE wallet_id = \$1`
		mock.
			ExpectQuery(query).
			WithArgs(walletID).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		tr := postgres.NewTransferRepository(db)

		_, err = tr.GetAllByWalletID(context.Background(), walletID)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrWalletNotExists {
			t.Errorf("unexpected error: %s, want: %s", err, customerrors.ErrWalletNotExists)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
