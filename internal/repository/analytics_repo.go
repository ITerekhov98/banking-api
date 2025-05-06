package repository

import (
	"context"
	"time"

	"banking-api/internal/db"
	"banking-api/internal/model"
)

type AnalyticsRepository struct{}

func (r *AnalyticsRepository) GetMonthlyStats(ctx context.Context, userID int64) (*model.MonthlyStats, error) {
	var income, expense float64

	incomeQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions t
		JOIN accounts a ON a.id = t.to_account_id
		WHERE a.user_id = $1 AND date_trunc('month', t.created_at) = date_trunc('month', CURRENT_DATE)
	`
	err := db.DB.QueryRowContext(ctx, incomeQuery, userID).Scan(&income)
	if err != nil {
		return nil, err
	}

	expenseQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM transactions t
		JOIN accounts a ON a.id = t.from_account_id
		WHERE a.user_id = $1 AND date_trunc('month', t.created_at) = date_trunc('month', CURRENT_DATE)
	`
	err = db.DB.QueryRowContext(ctx, expenseQuery, userID).Scan(&expense)
	if err != nil {
		return nil, err
	}

	return &model.MonthlyStats{
		Income:  income,
		Expense: expense,
	}, nil
}

func (r *AnalyticsRepository) GetUnpaidCreditTotal(ctx context.Context, userID int64) (float64, error) {
	var totalDebt float64

	query := `
		SELECT COALESCE(SUM(p_s.amount), 0)
		FROM payment_schedules p_s
		JOIN credits c ON c.id = p_s.credit_id
		WHERE p_s.paid = false AND c.user_id = $1
	`

	err := db.DB.QueryRowContext(ctx, query, userID).Scan(&totalDebt)
	return totalDebt, err
}

func (r *AnalyticsRepository) GetPredictedBalance(ctx context.Context, userID int64, until time.Time) (float64, error) {

	queryBalances := `
		SELECT SUM(balance::float)
		FROM accounts
		WHERE user_id = $1
	`
	var total float64
	err := db.DB.QueryRowContext(ctx, queryBalances, userID).Scan(&total)
	if err != nil {
		return 0, err
	}

	queryDeductions := `
		SELECT COALESCE(SUM(p_s.amount), 0)
		FROM payment_schedules p_s
		JOIN credits c ON c.id = p_s.credit_id
		WHERE p_s.paid = false AND p_s.due_date <= $1 AND c.user_id = $2
	`
	var due float64
	err = db.DB.QueryRowContext(ctx, queryDeductions, until, userID).Scan(&due)
	if err != nil {
		return 0, err
	}

	return total - due, nil
}
