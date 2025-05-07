package model

import "time"

// Credit represents a user-issued credit
// swagger:model Credit
type Credit struct {
	// ID of the credit
	// example: 101
	ID int64 `json:"id"`
	// ID of the user who owns the credit
	// example: 1
	UserID int64 `json:"user_id"`
	// ID of the associated account
	// example: 42
	AccountID int64 `json:"account_id"`
	// Original credit amount
	// example: 100000.00
	Principal float64 `json:"principal"`
	// Annual interest rate (percent)
	// example: 12.5
	InterestRate float64 `json:"interest_rate"`
	// Credit term in months
	// example: 12
	TermMonths int `json:"term_months"`
	// Monthly annuity payment
	// example: 8885.44
	MonthlyPayment float64 `json:"monthly_payment"`
	// Date the credit was issued
	// example: 2025-05-01T00:00:00Z
	CreatedAt time.Time `json:"created_at"`
}

// PaymentSchedule represents a monthly credit payment
// swagger:model PaymentSchedule
type PaymentSchedule struct {
	// Due date for the payment
	// example: 2025-06-01T00:00:00Z
	DueDate time.Time `json:"due_date"`
	// Payment amount for the month
	// example: 8885.44
	Amount float64 `json:"amount"`
	// Whether the payment was successfully made
	// example: false
	Paid bool `json:"paid"`
	// Whether a penalty was applied due to late or failed payment
	// example: true
	PenaltyApplied bool `json:"penalty_applied"`
}
