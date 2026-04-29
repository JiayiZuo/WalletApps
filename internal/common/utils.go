package common

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func CheckPassword(password, salt, storedHash string) bool {
	hash := sha256.Sum256([]byte(password + salt))
	calculatedHash := hex.EncodeToString(hash[:])
	return calculatedHash == storedHash
}

func MockGenerateAddress() (string, error) {
	// 以太坊地址是 20 字节 (160位)
	bytes := make([]byte, 20)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(bytes), nil
}
