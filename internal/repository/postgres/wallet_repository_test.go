package postgres_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/customerrors"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository/postgres"
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

	wr := postgres.NewWalletRepository(db)

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
		query := `SELECT wallet_id, balance FROM wallet WHERE wallet_id = \$1`
		mock.
			ExpectQuery(query).
			WithArgs(walletID).
			WillReturnRows(
				sqlmock.NewRows([]string{"wallet_id", "balance"}).AddRow(walletID, 100.0),
			)

		wr := postgres.NewWalletRepository(db)

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

		walletID := "NotExistingID"
		mock.
			ExpectQuery(`SELECT wallet_id, balance FROM wallet WHERE wallet_id = \$1`).
			WithArgs(walletID).
			WillReturnError(sql.ErrNoRows)

		wr := postgres.NewWalletRepository(db)

		_, err = wr.GetByWalletID(context.Background(), walletID)
		if err == nil {
			t.Errorf("expected error")
		}
		if err != customerrors.ErrWalletNotExists {
			t.Errorf("unexpected error: %s, got %s", customerrors.ErrWalletNotExists, err)
		}

		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
