package utils

import "math"

func CalculateAnnuity(principal float64, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 12 / 100
	return principal * monthlyRate / (1 - math.Pow(1+monthlyRate, float64(-months)))
}
