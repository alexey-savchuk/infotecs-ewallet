package postgres

import (
	"context"
	"database/sql"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{
		db: db,
	}
}

func (wr *WalletRepository) Create(ctx context.Context) (*repository.DBWallet, error) {
	query := `INSERT INTO wallet DEFAULT VALUES RETURNING wallet_id, balance`
	row := wr.db.QueryRowContext(ctx, query)
	wallet := &repository.DBWallet{}
	if err := row.Scan(&wallet.WalletID, &wallet.Balance); err != nil {
		return nil, err
	}
	return wallet, nil
}

func (w *WalletRepository) GetByWalletID(ctx context.Context, walletID string) (*repository.DBWallet, error) {
	query := `SELECT wallet_id, balance FROM wallet WHERE wallet_id = ?`
	row := w.db.QueryRowContext(ctx, query, walletID)
	wallet := repository.DBWallet{}
	if err := row.Scan(&wallet.WalletID, &wallet.Balance); err != nil {
		return nil, err
	}
	return &wallet, nil
}
