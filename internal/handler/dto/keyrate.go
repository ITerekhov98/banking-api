package dto

// KeyRateResponse contain current key rate
// swagger:model KeyRateResponse
type KeyRateResponse struct {
	// Current key rate
	// example: 21
	KeyRate float64 `json:"key_rate"`
}
