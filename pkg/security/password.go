package security

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 12
)

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
