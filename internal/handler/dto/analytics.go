package dto

// swagger:model predictedBalanceResponse
type PredictedBalanceResponse struct {
	// Predicted account balance
	// example: 123456.78
	PredictedBalance float64 `json:"predicted_balance"`
}
