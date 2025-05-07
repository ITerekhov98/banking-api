package dto

// createCreditRequest represents a credit creation request
// swagger:model createCreditRequest
type CreateCreditRequest struct {
	// ID of the account to associate the credit with
	// example: 42
	AccountID int64 `json:"account_id"`
	// Principal credit amount
	// example: 100000.00
	Principal float64 `json:"principal"`
	// Annual interest rate in percent
	// example: 12.5
	InterestRate float64 `json:"interest_rate"`
	// Term of the credit in months
	// example: 12
	TermMonths int `json:"term_months"`
}
