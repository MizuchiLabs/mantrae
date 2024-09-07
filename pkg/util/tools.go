// Package util contains various utility functions
package util

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenPassword generates a random password of the specified length
func GenPassword(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(bytes)[:length]
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// HashBasicAuth hashes a password using bcrypt (htpasswd)
func HashBasicAuth(userString string) (string, error) {
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

// IsValidURL checks if a URL is valid url string
func IsValidURL(u string) bool {
	// If no scheme is provided, prepend "http://"
	if !strings.Contains(u, "://") {
		u = "http://" + u
	}

	parsedURL, err := url.Parse(u)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return false
	}

	host := parsedURL.Hostname()
	port := parsedURL.Port()
	if port != "" {
		if _, err = net.LookupPort("tcp", port); err != nil {
			return false
		}
	}

	// Check if it's an IP address (including loopback)
	ip := net.ParseIP(host)
	if ip != nil {
		return true
	}

	// Check if it's localhost
	if !strings.Contains(host, ".") {
		_, err = strconv.Atoi(host)
		return err != nil // Valid if it's not just a number
	}

	// Check if it's a valid domain name
	domainRegex := `^([a-zA-Z0-9-]+\.)*[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(domainRegex, host)
	if err != nil {
		return false
	}
	return matched
}

// IsValidEmail checks if an email is valid
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return false
	}
	return matched
}
