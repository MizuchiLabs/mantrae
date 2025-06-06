package util

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// ExtractBaseDomain returns the base domain of a URL
func ExtractBaseDomain(domain string) (string, error) {
	// Ensure the domain doesn't contain a scheme
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	return publicsuffix.EffectiveTLDPlusOne(domain)
}

// ExtractDomainFromRule extracts all the domains from a traefik rule (e.g. Host(`domain.com`))
func ExtractDomainFromRule(rule string) ([]string, error) {
	re := regexp.MustCompile(`Host\(` + "`" + `([^` + "`" + `]+)` + "`" + `\)`)
	matches := re.FindAllStringSubmatch(rule, -1)

	if len(matches) == 0 {
		return nil, fmt.Errorf("no domains found in rule")
	}

	var domains []string
	for _, match := range matches {
		if len(match) > 1 {
			domains = append(domains, match[1])
		}
	}
	return domains, nil
}
