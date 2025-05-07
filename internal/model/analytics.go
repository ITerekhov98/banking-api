package model

// MonthlyStats represents a summary of income, expenses and credit obligations
// swagger:model MonthlyStats
type MonthlyStats struct {
	// Total income for the current month
	// example: 120000.00
	Income float64 `json:"monthly_income"`
	// Total expenses for the current month
	// example: 80000.00
	Expense float64 `json:"monthly_expense"`
	// Total unpaid credit obligations
	// example: 50000.00
	DueTotal float64 `json:"credit_due_total"`
}
