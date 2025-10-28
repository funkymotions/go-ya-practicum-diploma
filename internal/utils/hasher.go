package utils

import (
	"crypto/md5"
	"fmt"
)

// using one-way md5 hash algorithm for educational purposes only
func HashPassword(password string) string {
	hasher := md5.New()
	hasher.Write([]byte(password))
	hash := hasher.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func ComparePasswords(hashedPassword, password string) bool {
	return hashedPassword == HashPassword(password)
}
