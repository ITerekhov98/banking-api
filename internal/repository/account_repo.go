package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"banking-api/internal/db"
	"banking-api/internal/model"
)

type AccountRepository struct{}

func (r *AccountRepository) Create(ctx context.Context, userID int64) (*model.Account, error) {
	query := `INSERT INTO accounts (user_id, balance) VALUES ($1, 0) RETURNING id, balance`
	row := db.DB.QueryRowContext(ctx, query, userID)

	var acc model.Account
	var balanceStr string

	acc.UserID = userID
	err := row.Scan(&acc.ID, &balanceStr)
	if err != nil {
		return nil, err
	}
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return nil, err
	}
	acc.Balance = balance
	return &acc, nil
}

func (r *AccountRepository) TransferTx(ctx context.Context, fromID, toID int64, amount float64) error {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check sender balance
	var balanceStr string
	err = tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, fromID).Scan(&balanceStr)
	if err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}

	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return err
	}
	if balance < amount {
		return errors.New("insufficient funds")
	}

	// withdraw
	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		return err
	}

	// deposit
	_, err = tx.ExecContext(ctx, `UPDATE accounts SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		return err
	}

	// log transaction
	_, err = tx.ExecContext(ctx, `
		INSERT INTO transactions (from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3)
	`, fromID, toID, amount)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *AccountRepository) GetOwnerID(ctx context.Context, accountID int64) (int64, error) {
	query := `SELECT user_id FROM accounts WHERE id = $1`
	var ownerID int64
	err := db.DB.QueryRowContext(ctx, query, accountID).Scan(&ownerID)
	return ownerID, err
}
