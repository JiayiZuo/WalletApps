package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func CheckPassword(password, salt, storedHash string) bool {
	hash := sha256.Sum256([]byte(password + salt))
	calculatedHash := hex.EncodeToString(hash[:])
	return calculatedHash == storedHash
}
