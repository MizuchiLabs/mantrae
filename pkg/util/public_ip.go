package util

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// Public IP APIs
var ipAPIs = []string{
	"https://api.ipify.org?format=text",
	"https://ifconfig.co/ip",
	"https://checkip.amazonaws.com",
	"https://ipinfo.io/ip",
}

func getIP(url string) (string, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("non-200 response from API")
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

func isValidPublicIP(ip string) bool {
	if ip == "" {
		return false
	}

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	// Check if it's a private or loopback IP
	if parsedIP.IsLoopback() || parsedIP.IsPrivate() {
		return false
	}

	return true
}

func GetPublicIP() (string, error) {
	for _, api := range ipAPIs {
		ip, err := getIP(api)
		if err == nil && isValidPublicIP(ip) {
			return ip, nil
		}
		slog.Warn("Failed to query API", "API", api, "Error", err)
	}
	return "", fmt.Errorf("failed to get public IP")
}
