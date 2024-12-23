package utils

import (
	"crypto/sha512"
	"encoding/hex"
)

const saltSize = 16

func HashPassword(password string, salt string) string {
	var sha512Hasher = sha512.New()
	sha512Hasher.Write(append([]byte(password), salt...))
	var hashedPasswordBytes = sha512Hasher.Sum(nil)
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)
	return hashedPasswordHex
}

func HashWithSalt(password string) (string, string) {
	var salt = RandomString(saltSize)
	hashed := HashPassword(password, salt)
	return hashed, salt
}

// Check if two passwords match
func CheckPassword(
	password string,
	hashedPassword,
	salt string) bool {
	if len(password) == 0 {
		return false
	}
	var currPasswordHash = HashPassword(password, salt)
	return hashedPassword == currPasswordHash
}
