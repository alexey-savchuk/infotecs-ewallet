package postgres

import (
	"context"
	"database/sql"

	"github.com/alexey-savchuk/infotecs-ewallet/internal/customerrors"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/repository"
	"github.com/alexey-savchuk/infotecs-ewallet/internal/service"
	"github.com/shopspring/decimal"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{
		db: db,
	}
}

func (tr *TransferRepository) Create(ctx context.Context, transfer *service.Transfer) (*repository.DBTransfer, error) {
	tx, err := tr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	query := `SELECT balance FROM wallet WHERE wallet_id = ?`
	row := tx.QueryRowContext(ctx, query, transfer.FromWallet)
	var balance decimal.Decimal
	err = row.Scan(&balance)
	if err != nil {
		return nil, customerrors.ErrFromWalletNotExists
	}
	if balance.LessThan(transfer.Amount) {
		return nil, customerrors.ErrInsufficientFunds
	}

	query = `UPDATE wallet SET balance = balance - ? WHERE wallet_id = ?`
	_, err = tx.ExecContext(ctx, query, transfer.Amount, transfer.FromWallet)
	if err != nil {
		return nil, customerrors.ErrFromWalletNotExists
	}
	query = `UPDATE wallet SET balance = balance + ? WHERE wallet_id = ?`
	_, err = tx.ExecContext(ctx, query, transfer.Amount, transfer.ToWallet)
	if err != nil {
		return nil, customerrors.ErrToWalletNotExists
	}

	query = `INSERT INTO transfer (from_wallet, to_wallet, amount)
			 VALUES (?, ?, ?)
			 RETURNING transfer_id, time, from_wallet, to_wallet, amount`
	row = tx.QueryRowContext(
		ctx, query, transfer.FromWallet, transfer.ToWallet, transfer.Amount,
	)
	var dbTransfer repository.DBTransfer
	err = row.Scan(&dbTransfer.TransferID, &dbTransfer.Time, &dbTransfer.FromWallet, &dbTransfer.ToWallet, &dbTransfer.Amount)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &dbTransfer, nil
}

func (tr *TransferRepository) GetAllByWalletID(ctx context.Context, walletID string) ([]repository.DBTransfer, error) {
	tx, err := tr.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable, ReadOnly: true})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() //nolint:errcheck

	query := `SELECT wallet_id FROM wallet WHERE wallet_id = ?`
	row := tx.QueryRowContext(ctx, query, walletID)
	err = row.Scan(&walletID)
	if err != nil {
		return nil, customerrors.ErrWalletNotExists
	}

	query = `SELECT transfer_id, time, from_wallet, to_wallet, amount
			  FROM transfer WHERE from_wallet = ? OR to_wallet = ? ORDER BY 2 DESC`
	rows, err := tr.db.QueryContext(ctx, query, walletID, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dbTransfers := make([]repository.DBTransfer, 0)
	for rows.Next() {
		var dbTransfer repository.DBTransfer
		err = rows.Scan(&dbTransfer.TransferID, &dbTransfer.Time, &dbTransfer.FromWallet, &dbTransfer.ToWallet, &dbTransfer.Amount)
		if err != nil {
			return nil, err
		}
		dbTransfers = append(dbTransfers, dbTransfer)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return dbTransfers, nil
}
