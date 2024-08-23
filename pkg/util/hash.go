package util

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(userString string) (string, error) {
	if userString == "" {
		return "", nil
	}

	user := strings.Split(userString, ":")[0]
	password := strings.Split(userString, ":")[1]

	// If the password is already hashed, return it
	if strings.HasPrefix(password, "$2") && len(password) >= 60 {
		return user + ":" + password, nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return user + ":" + string(hash), nil
}
