package util

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// IPAddresses holds both IPv4 and IPv6 addresses
type IPAddresses struct {
	IPv4 string `json:"ipv4,omitempty"`
	IPv6 string `json:"ipv6,omitempty"`
}

var (
	cachedIPs     IPAddresses
	cacheExpiry   time.Time
	cacheDuration = 5 * time.Minute
)

func GetPublicIPsCached() (IPAddresses, error) {
	if time.Now().Before(cacheExpiry) {
		return cachedIPs, nil
	}

	ips, err := GetPublicIPs()
	if err != nil {
		return IPAddresses{}, err
	}

	// Update cache
	cachedIPs = ips
	cacheExpiry = time.Now().Add(cacheDuration)

	return ips, nil
}

func GetPublicIPs() (IPAddresses, error) {
	ipv4Services := []string{
		"https://api.ipify.org",
		"https://api.ip.sb/ip",
		"https://ipinfo.io/ip",
		"https://ifconfig.me/ip",
		"https://ifconfig.co/ip",
		"https://ipecho.net/plain",
		"https://checkip.amazonaws.com",
	}

	ipv6Services := []string{
		"https://api6.ipify.org",
		"https://v6.ident.me",
		"https://api-ipv6.ip.sb/ip",
	}

	// Create channels for concurrent IP fetching
	ipv4Ch := make(chan string, 1)
	ipv6Ch := make(chan string, 1)
	errorCh := make(chan error, 2)

	// Fetch IPv4 and IPv6 addresses concurrently
	go func() {
		ip, err := getIP(ipv4Services, IsValidIPv4)
		if err != nil {
			errorCh <- fmt.Errorf("IPv4: %v", err)
			ipv4Ch <- ""
			return
		}
		ipv4Ch <- ip
	}()

	go func() {
		ip, err := getIP(ipv6Services, IsValidIPv6)
		if err != nil {
			errorCh <- fmt.Errorf("IPv6: %v", err)
			ipv6Ch <- ""
			return
		}
		ipv6Ch <- ip
	}()

	// Wait for results
	ipv4 := <-ipv4Ch
	ipv6 := <-ipv6Ch

	// Check if we got at least one address
	if ipv4 == "" && ipv6 == "" {
		return IPAddresses{}, fmt.Errorf("failed to get any IP address")
	}

	return IPAddresses{
		IPv4: ipv4,
		IPv6: ipv6,
	}, nil
}

func getIP(services []string, validationFunc func(string) bool) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	for _, service := range services {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", service, nil)
		if err != nil {
			continue
		}
		req.Header.Set(
			"User-Agent",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
		)

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}
		ip := strings.TrimSpace(string(body))
		if validationFunc(ip) {
			return ip, nil
		}
	}

	return "", fmt.Errorf("failed to get IP from any service")
}

func IsValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.To4() != nil
}

func IsValidIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}
	return parsedIP.To4() == nil && parsedIP.To16() != nil
}

func GetLocalIP() string {
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "3000"
	}
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return fmt.Sprintf("http://%s:%s", ipnet.IP.String(), port)
			}
		}
	}
	return fmt.Sprintf("http://127.0.0.1:%s", port)
}
