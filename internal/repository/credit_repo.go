package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"banking-api/internal/db"
	"banking-api/internal/model"
	"banking-api/pkg/logger"
)

const PaymentPenalty float64 = 0.10

type CreditRepository struct{}

func (r *CreditRepository) SaveCredit(ctx context.Context, c *model.Credit) error {
	query := `
		INSERT INTO credits (user_id, account_id, principal, interest_rate, term_months, monthly_payment)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at;
	`
	return db.DB.QueryRowContext(ctx, query, c.UserID, c.AccountID, c.Principal, c.InterestRate, c.TermMonths, c.MonthlyPayment).
		Scan(&c.ID, &c.CreatedAt)
}

func (r *CreditRepository) GenerateSchedule(ctx context.Context, creditID int64, months int, amount float64) error {
	query := `INSERT INTO payment_schedules (credit_id, due_date, amount) VALUES ($1, $2, $3)`
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	start := time.Now()
	for i := 1; i <= months; i++ {
		due := start.AddDate(0, i, 0)
		_, err := tx.ExecContext(ctx, query, creditID, due, amount)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *CreditRepository) GetScheduleByCreditID(ctx context.Context, userID, creditID int64) ([]model.PaymentSchedule, error) {
	// Check owner
	var ownerID int64
	err := db.DB.QueryRowContext(ctx, `SELECT user_id FROM credits WHERE id = $1`, creditID).Scan(&ownerID)
	if err != nil {
		return nil, fmt.Errorf("credit not found")
	}
	if ownerID != userID {
		return nil, fmt.Errorf("access denied")
	}

	query := `
		SELECT due_date, amount, paid, penalty_applied
		FROM payment_schedules
		WHERE credit_id = $1
		ORDER BY due_date
	`

	rows, err := db.DB.QueryContext(ctx, query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.PaymentSchedule
	for rows.Next() {
		var p model.PaymentSchedule
		err := rows.Scan(&p.DueDate, &p.Amount, &p.Paid, &p.PenaltyApplied)
		if err != nil {
			return nil, err
		}
		result = append(result, p)
	}

	return result, nil
}

func (r *CreditRepository) CheckDuePayments(ctx context.Context, now time.Time) error {
	query := `
		SELECT ps.id, ps.credit_id, ps.amount, c.account_id
		FROM payment_schedules ps
		JOIN credits c ON c.id = ps.credit_id
		WHERE ps.due_date <= $1 AND ps.paid = false
	`
	rows, err := db.DB.QueryContext(ctx, query, now)
	if err != nil {
		return err
	}
	defer rows.Close()

	type payment struct {
		ID        int64
		CreditID  int64
		AccountID int64
		Amount    float64
	}

	var payments []payment
	for rows.Next() {
		var p payment
		if err := rows.Scan(&p.ID, &p.CreditID, &p.Amount, &p.AccountID); err != nil {
			return err
		}
		payments = append(payments, p)
	}

	for _, p := range payments {
		err := r.processPayment(ctx, p)
		if err != nil {
			logger.Error("payment failed:", err)
		}
	}

	return nil
}

func (r *CreditRepository) processPayment(ctx context.Context, p struct {
	ID, CreditID, AccountID int64
	Amount                  float64
}) error {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balanceStr string
	err = tx.QueryRowContext(ctx, `SELECT balance FROM accounts WHERE id = $1 FOR UPDATE`, p.AccountID).Scan(&balanceStr)
	if err != nil {
		return err
	}
	balance, err := strconv.ParseFloat(balanceStr, 64)
	if err != nil {
		return err
	}
	if balance >= p.Amount {
		_, err := tx.ExecContext(ctx, `
			UPDATE accounts SET balance = balance - $1 WHERE id = $2
		`, p.Amount, p.AccountID)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE payment_schedules SET paid = true WHERE id = $1
		`, p.ID)
		if err != nil {
			return err
		}
	} else {
		// Начисляем штраф
		penalty := p.Amount * PaymentPenalty
		newAmount := p.Amount + penalty

		_, err = tx.ExecContext(ctx, `
			UPDATE payment_schedules SET amount = $1, penalty_applied = true WHERE id = $2
		`, newAmount, p.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
