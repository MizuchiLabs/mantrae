// Package util contains various utility functions
package util

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// IsTest returns true if the current program is running in a test environment
func IsTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

func SafeDeref(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

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

// IsHtpasswdFormat checks if a string is already in htpasswd format
func IsHtpasswdFormat(s string) bool {
	// HTPasswd formats we support:
	// - bcrypt: $2y$ or $2a$ or $2b$ followed by cost and 22 chars
	// - MD5: $apr1$ followed by salt and hash
	bcryptRegex := regexp.MustCompile(`^\$2[ayb]\$.{56}$`)
	md5Regex := regexp.MustCompile(`^\$apr1\$.{30,}$`)

	return bcryptRegex.MatchString(s) || md5Regex.MatchString(s)
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

func ValidSSLCert(domain string) error {
	conn, err := tls.DialWithDialer(
		&net.Dialer{Timeout: 5 * time.Second},
		"tcp",
		domain+":443",
		&tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return fmt.Errorf("timeout reached after 5 seconds: %v", err)
		}
		return fmt.Errorf("could not establish TLS connection: %v", err)
	}
	defer conn.Close()

	// Get the certificate
	cert := conn.ConnectionState().PeerCertificates[0]

	// Check if the certificate is currently valid
	now := time.Now()
	if now.Before(cert.NotBefore) {
		return fmt.Errorf("certificate is not yet valid")
	}
	if now.After(cert.NotAfter) {
		return fmt.Errorf("certificate has expired")
	}

	// Ensure the domain has a scheme
	if !strings.HasPrefix(domain, "http://") && !strings.HasPrefix(domain, "https://") {
		domain = "https://" + domain
	}

	// Extra cloudflare checks
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make a request to the domain
	resp, err := client.Get(domain)
	if err != nil {
		return fmt.Errorf("error reaching domain: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP status codes commonly returned by Cloudflare for SSL issues
	switch resp.StatusCode {
	case 500:
		return fmt.Errorf("Internal Server Error (Error 500)")
	case 502:
	case 504:
		return fmt.Errorf("Bad Gateway (Error 502 or 504)")
	case 503:
		return fmt.Errorf("Service Unavailable (Error 503)")
	case 520:
		return fmt.Errorf("Web server returns an unknown error (Error 520)")
	case 521:
		return fmt.Errorf("Web server is down or unreachable (Error 521)")
	case 522:
		return fmt.Errorf("Connection timed out (Error 522)")
	case 523:
		return fmt.Errorf("Origin is not reachable (Error 523)")
	case 524:
		return fmt.Errorf("A timeout occurred during SSL handshake (Error 524)")
	case 525:
		return fmt.Errorf("SSL handshake failed between Cloudflare and origin (Error 525)")
	case 526:
		return fmt.Errorf("Invalid SSL certificate on the origin server (Error 526)")
	}

	return nil
}
