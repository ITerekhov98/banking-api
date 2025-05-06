package security

import (
	"errors"
	"os"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHMAC(data string) (string, error) {
	secret := os.Getenv("HMAC_SECRET")
	if secret == "" {
		return "", errors.New("hmac secret nor found")
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil)), nil
}
