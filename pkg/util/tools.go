// Package util contains various utility functions
package util

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	Username  string `json:"username,omitempty"`
	ServerURL string `json:"serverUrl,omitempty"`
	ProfileID int64  `json:"profileId,omitempty"`
	Secret    string `json:"secret,omitempty"`
	jwt.RegisteredClaims
}

// IsTest returns true if the current program is running in a test environment
func IsTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}

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

// EncodeUserJWT generates a JWT for user login
func EncodeUserJWT(username string) (string, error) {
	if username == "" {
		return "", errors.New("username cannot be empty")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(App.Secret))
}

// EncodeAgentJWT generates a JWT for the agent
func EncodeAgentJWT(profileID int64) (string, error) {
	if profileID == 0 {
		return "", errors.New("profileID cannot be empty")
	}

	expirationTime := time.Now().Add(14 * 24 * time.Hour) // 14 days
	claims := Claims{
		ServerURL: App.ServerURL + ":" + App.AgentPort,
		ProfileID: profileID,
		Secret:    App.Secret,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(App.Secret))
}

// DecodeJWT decodes the token and returns claims if valid
func DecodeJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(App.Secret), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
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

// ExtractDomainFromRule extracts the domain inside a Host(`domain.com`) rule
func ExtractDomainFromRule(rule string) (string, error) {
	re := regexp.MustCompile(`Host\(` + "`" + `([^` + "`" + `]+)` + "`" + `\)`)
	matches := re.FindStringSubmatch(rule)
	if len(matches) < 2 {
		return "", fmt.Errorf("no domain found in rule")
	}
	return matches[1], nil
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

func IsRunningInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}
