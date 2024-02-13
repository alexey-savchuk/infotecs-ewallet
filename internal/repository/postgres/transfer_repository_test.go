package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/service"
	"github.com/shopspring/decimal"
)

func TestTransferRepository_Create(t *testing.T) {
	t.Run("WalletsExist", func(t *testing.T) {
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
		query := `UPDATE wallet SET balance = balance \- \? WHERE wallet_id = \?`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.FromWallet).
			WillReturnResult(sqlmock.NewResult(1, 1))
		query = `UPDATE wallet SET balance = balance \+ \? WHERE wallet_id = \?`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.ToWallet).
			WillReturnResult(sqlmock.NewResult(1, 1))
		query = `INSERT INTO transfer \(from_wallet, to_wallet, amount\)
				 VALUES \(\?, \?, \?\)
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

		tr := NewTransferRepository(db)

		_, err = tr.Create(context.Background(), transfer)
		if err != nil {
			t.Errorf("unexpected error: %s", err)
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
			FromWallet: "1",
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}

		mock.ExpectBegin()
		query := `UPDATE wallet SET balance = balance \- \? WHERE wallet_id = \?`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.FromWallet).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		tr := NewTransferRepository(db)

		_, err = tr.Create(context.Background(), transfer)
		if err == nil {
			t.Errorf("expected error")
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
			ToWallet:   "2",
			Amount:     decimal.NewFromFloat(100.0),
		}

		mock.ExpectBegin()
		query := `UPDATE wallet SET balance = balance \- \? WHERE wallet_id = \?`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.FromWallet).
			WillReturnResult(sqlmock.NewResult(1, 1))
		query = `UPDATE wallet SET balance = balance \+ \? WHERE wallet_id = \?`
		mock.
			ExpectExec(query).
			WithArgs(transfer.Amount, transfer.ToWallet).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		tr := NewTransferRepository(db)

		_, err = tr.Create(context.Background(), transfer)
		if err == nil {
			t.Errorf("expected error")
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
		query := `SELECT transfer_id, time, from_wallet, to_wallet, amount
				  FROM transfer WHERE from_wallet = \? OR to_wallet = \? ORDER BY 2 DESC`
		mock.
			ExpectQuery(query).
			WithArgs(walletID, walletID).
			WillReturnRows(
				sqlmock.
					NewRows([]string{"transfer_id", "time", "from_wallet", "to_wallet", "amount"}).
					AddRow("1", time.Now(), walletID, "2", 100.0).
					AddRow("2", time.Now(), "3", walletID, 200.0),
			)

		tr := NewTransferRepository(db)

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

		walletID := "1"
		query := `SELECT transfer_id, time, from_wallet, to_wallet, amount
				  FROM transfer WHERE from_wallet = \? OR to_wallet = \? ORDER BY 2 DESC`
		mock.
			ExpectQuery(query).
			WithArgs(walletID, walletID).
			WillReturnError(sql.ErrNoRows)

		tr := NewTransferRepository(db)

		_, err = tr.GetAllByWalletID(context.Background(), walletID)
		if err == nil {
			t.Errorf("expected error")
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
