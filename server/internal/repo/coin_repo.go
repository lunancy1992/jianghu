package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lunancy1992/jianghu-server/internal/model"
)

type CoinRepo struct {
	db *sql.DB
}

func NewCoinRepo(db *sql.DB) *CoinRepo {
	return &CoinRepo{db: db}
}

func (r *CoinRepo) GetAccount(ctx context.Context, userID int64) (*model.CoinAccount, error) {
	a := &model.CoinAccount{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, balance, created_at, updated_at FROM coin_accounts WHERE user_id=?`,
		userID).Scan(&a.ID, &a.UserID, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *CoinRepo) CreateAccount(ctx context.Context, userID int64, initialBalance int64) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`INSERT OR IGNORE INTO coin_accounts (user_id, balance, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		userID, initialBalance, now, now,
	)
	return err
}

// Deduct subtracts coins. Returns error if insufficient balance. Should be called within a tx.
func (r *CoinRepo) Deduct(ctx context.Context, tx *sql.Tx, userID int64, amount int64, txType, refID, note string) error {
	var balance int64
	err := tx.QueryRowContext(ctx,
		`SELECT balance FROM coin_accounts WHERE user_id=?`, userID).Scan(&balance)
	if err != nil {
		return fmt.Errorf("get balance: %w", err)
	}
	if balance < amount {
		return fmt.Errorf("insufficient balance: have %d, need %d", balance, amount)
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE coin_accounts SET balance = balance - ?, updated_at = ? WHERE user_id = ?`,
		amount, time.Now(), userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO coin_transactions (user_id, amount, type, ref_id, note, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		userID, -amount, txType, refID, note, time.Now())
	return err
}

// Credit adds coins. Should be called within a tx.
func (r *CoinRepo) Credit(ctx context.Context, tx *sql.Tx, userID int64, amount int64, txType, refID, note string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE coin_accounts SET balance = balance + ?, updated_at = ? WHERE user_id = ?`,
		amount, time.Now(), userID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO coin_transactions (user_id, amount, type, ref_id, note, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		userID, amount, txType, refID, note, time.Now())
	return err
}

func (r *CoinRepo) ListTransactions(ctx context.Context, userID int64, page, size int) ([]*model.CoinTransaction, int, error) {
	var total int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM coin_transactions WHERE user_id=?`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, amount, type, ref_id, note, created_at
		 FROM coin_transactions WHERE user_id=?
		 ORDER BY created_at DESC LIMIT ? OFFSET ?`, userID, size, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*model.CoinTransaction
	for rows.Next() {
		t := &model.CoinTransaction{}
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Type, &t.RefID, &t.Note, &t.CreatedAt); err != nil {
			return nil, 0, err
		}
		list = append(list, t)
	}
	return list, total, rows.Err()
}

// GetActiveUsers returns user IDs who have been active (commented or liked) recently.
func (r *CoinRepo) GetActiveUsers(ctx context.Context, since time.Time) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT user_id FROM comments WHERE created_at >= ?
		 UNION
		 SELECT DISTINCT user_id FROM likes WHERE created_at >= ?`,
		since, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// BeginTx starts a new database transaction.
func (r *CoinRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
