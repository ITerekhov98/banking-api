package model

import "time"

type Credit struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	AccountID      int64     `json:"account_id"`
	Principal      float64   `json:"principal"`
	InterestRate   float64   `json:"interest_rate"`
	TermMonths     int       `json:"term_months"`
	MonthlyPayment float64   `json:"monthly_payment"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaymentSchedule struct {
	DueDate        time.Time `json:"due_date"`
	Amount         float64   `json:"amount"`
	Paid           bool      `json:"paid"`
	PenaltyApplied bool      `json:"penalty_applied"`
}
