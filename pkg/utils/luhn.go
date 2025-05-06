package utils

import (
	"math/rand"
	"strconv"
)

// generate card number via luhn algorithm
func GenerateCardNumber(prefix string) string {

	number := prefix
	for range 15 - len(prefix) {
		number += strconv.Itoa(rand.Intn(10))
	}

	sum := 0
	for i := range 15 {
		n, _ := strconv.Atoi(string(number[14-i]))
		if i%2 == 0 {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
	}
	checkDigit := (10 - sum%10) % 10
	return number + strconv.Itoa(checkDigit)

}
