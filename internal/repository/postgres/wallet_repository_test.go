package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestWalletRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	query := `INSERT INTO wallet DEFAULT VALUES RETURNING wallet_id, balance`
	mock.
		ExpectQuery(query).
		WillReturnRows(
			sqlmock.NewRows([]string{"wallet_id", "balance"}).AddRow("1", 100.0),
		)

	wr := NewWalletRepository(db)

	_, err = wr.Create(context.Background())
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestWalletRepository_GetByWalletID(t *testing.T) {
	t.Run("WalletExists", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		walletID := "1"
		query := `SELECT wallet_id, balance FROM wallet WHERE wallet_id = \?`
		mock.
			ExpectQuery(query).
			WithArgs(walletID).
			WillReturnRows(
				sqlmock.NewRows([]string{"wallet_id", "balance"}).AddRow(walletID, 100.0),
			)

		wr := NewWalletRepository(db)

		_, err = wr.GetByWalletID(context.Background(), walletID)
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
		mock.
			ExpectQuery(`SELECT wallet_id, balance FROM wallet WHERE wallet_id = \?`).
			WithArgs(walletID).
			WillReturnError(sql.ErrNoRows)

		wr := NewWalletRepository(db)

		_, err = wr.GetByWalletID(context.Background(), walletID)
		if err == nil {
			t.Errorf("expected error")
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
