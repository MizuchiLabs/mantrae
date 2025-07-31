package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenPassword generates a random password of the specified length
func GenPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))
	for i := range password {
		index, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return ""
		}
		password[i] = charset[index.Int64()]
	}
	return string(password)
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsHtpasswdFormat checks if a string is already in htpasswd format
func IsHtpasswdFormat(s string) bool {
	// HTPasswd formats we support:
	// - bcrypt: $2y$ or $2a$ or $2b$ followed by cost and 22 chars
	// - MD5: $apr1$ followed by salt and hash
	bcryptRegex := regexp.MustCompile(`^\$2[ayb]\$.{56}$`)
	md5Regex := regexp.MustCompile(`^\$apr1\$.{30,}$`)

	return bcryptRegex.MatchString(s) || md5Regex.MatchString(s)
}

// GenerateToken creates a random url safe token of the specified length
func GenerateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	token := base32.StdEncoding.EncodeToString(b)
	return strings.ToLower(strings.TrimRight(token, "="))
}

func GenerateAgentToken(profileID, agentID string) string {
	return fmt.Sprintf("%s.%s.%s", profileID, agentID, GenerateToken(8))
}

// GenerateOTP creates a secure 6-digit token
func GenerateOTP() (string, error) {
	const otpChars = "0123456789"
	buffer := make([]byte, 6)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otp := []rune("0123456789")
	otpLength := len(otpChars)
	token := make([]rune, 6)

	for i := range buffer {
		token[i] = otp[int(buffer[i])%otpLength]
	}

	return string(token), nil
}

func HashOTP(otp string) string {
	sum := sha256.Sum256([]byte(otp))
	return hex.EncodeToString(sum[:])
}
