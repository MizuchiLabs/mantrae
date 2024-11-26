package client

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"golang.org/x/exp/slices"
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

func GetPrivateIP() ([]string, error) {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	excluded := []string{"lo", "docker", "br-", "veth", "kube", "cni"}
	for _, iface := range interfaces {
		if slices.ContainsFunc(excluded, func(s string) bool {
			return strings.Contains(iface.Name, s)
		}) || iface.Flags&net.FlagUp == 0 {
			continue
		} else {
			addrs, err := iface.Addrs()
			if err != nil {
				return nil, err
			}

			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ips = append(ips, ipnet.IP.String())
					}
				}
			}
		}
	}

	if len(ips) == 0 {
		return nil, errors.New("no private IP addresses found")
	}

	return ips, nil
}
