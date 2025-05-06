package model

type MonthlyStats struct {
	Income   float64 `json:"monthly_income"`
	Expense  float64 `json:"monthly_expense"`
	DueTotal float64 `json:"credit_due_total"`
}
