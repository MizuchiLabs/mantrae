package util

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
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

func isVirtualInterface(name string) bool {
	skipIfacesPrefixes := []string{
		"docker",
		"br-",
		"veth",
		"lo",
		"tun",
		"tap",
		"vmnet",
		"wg",
		"utun",
		"ppp",
		"virbr",
		"cni",
		"flannel",
		"weave",
	}
	name = strings.ToLower(name)
	for _, prefix := range skipIfacesPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}
	return false
}

func GetHostIPv4() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback
		}
		if isVirtualInterface(iface.Name) {
			continue // skip virtual interfaces
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip4 := ip.To4()
			if ip4 == nil {
				continue
			}
			// Got a good candidate
			return ip4.String(), nil
		}
	}

	return "", errors.New("no suitable IPv4 address found")
}

func GetPublicIPs() (IPAddresses, error) {
	if time.Now().Before(cacheExpiry) {
		return cachedIPs, nil
	}

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

	ips := IPAddresses{IPv4: ipv4, IPv6: ipv6}

	// Update cache
	cachedIPs = ips
	cacheExpiry = time.Now().Add(cacheDuration)

	return ips, nil
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
		defer func() {
			if err = resp.Body.Close(); err != nil {
				slog.Error("failed to close response body", "error", err)
			}
		}()

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

// GetClientIP returns the IP address of the client
func GetClientIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// ParseHostPort parses a URL string and returns the host and port.
func ParseHostPort(rawURL string) (host, port string) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", ""
	}

	h := u.Host
	if strings.Contains(h, ":") {
		host, port, err = net.SplitHostPort(h)
		if err != nil {
			return h, ""
		}
		return host, port
	}
	return h, ""
}

// OriginOnly returns the origin of a URL, without the path, query, or fragment.
func OriginOnly(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	u.Path, u.RawQuery, u.Fragment = "", "", ""
	return strings.TrimSuffix(u.String(), "/")
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

func CleanURL(url string) string {
	url = strings.TrimSpace(url)
	if !strings.HasPrefix(url, "http://") &&
		!strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	return url
}
