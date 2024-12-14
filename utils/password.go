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
	hashedPassword,
	currPassword string,
	salt string) bool {
	var currPasswordHash = HashPassword(currPassword, salt)
	if len(currPassword) == 0 {
		return false
	}
	return hashedPassword == currPasswordHash
}
