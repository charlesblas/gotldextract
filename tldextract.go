package gotldextract

import (
	"fmt"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// Result represents the extracted parts of a domain
type Result struct {
	Subdomain string
	Domain    string
	TLD       string
}

// Extract parses a domain/URL and extracts its parts
func Extract(domain string) (*Result, error) {
	// Clean the domain (remove protocol if present)
	domain = cleanDomain(domain)
	
	// Get the public suffix (TLD)
	suffix, icann := publicsuffix.PublicSuffix(domain)
	if !icann {
		// If not an ICANN suffix, treat as private
		suffix, _ = publicsuffix.PublicSuffix(domain)
	}
	
	// If the whole domain is just the suffix, it's not a valid domain
	if domain == suffix {
		return &Result{
			Subdomain: "",
			Domain:    "",
			TLD:       suffix,
		}, nil
	}
	
	// Remove the suffix to get the rest
	domainWithoutSuffix := strings.TrimSuffix(domain, "."+suffix)
	
	// Split by dots
	parts := strings.Split(domainWithoutSuffix, ".")
	
	if len(parts) == 0 {
		return &Result{
			Subdomain: "",
			Domain:    "",
			TLD:       suffix,
		}, nil
	}
	
	// The last part is the domain name
	domainName := parts[len(parts)-1]
	
	// Everything before is subdomain
	subdomain := ""
	if len(parts) > 1 {
		subdomain = strings.Join(parts[:len(parts)-1], ".")
	}
	
	return &Result{
		Subdomain: subdomain,
		Domain:    domainName,
		TLD:       suffix,
	}, nil
}

// ExtractFromURL is an alias for Extract that emphasizes it can handle URLs
func ExtractFromURL(url string) (*Result, error) {
	return Extract(url)
}

// Update updates the public suffix list
// Note: The golang.org/x/net/publicsuffix package uses an embedded list
// that is updated when the package itself is updated.
func Update() error {
	return fmt.Errorf("update not supported: the public suffix list is embedded in the golang.org/x/net/publicsuffix package")
}

// cleanDomain removes protocol and path from a URL to get just the domain
func cleanDomain(domain string) string {
	// Remove protocol
	if strings.Contains(domain, "://") {
		parts := strings.SplitN(domain, "://", 2)
		if len(parts) == 2 {
			domain = parts[1]
		}
	}
	
	// Remove path
	if idx := strings.Index(domain, "/"); idx != -1 {
		domain = domain[:idx]
	}
	
	// Remove port
	if idx := strings.Index(domain, ":"); idx != -1 {
		domain = domain[:idx]
	}
	
	// Remove trailing dot
	domain = strings.TrimSuffix(domain, ".")
	
	return strings.ToLower(domain)
}

// String returns the full domain (domain + TLD)
func (r *Result) String() string {
	if r.Domain == "" {
		return r.TLD
	}
	return r.Domain + "." + r.TLD
}

// FQDN returns the fully qualified domain name
func (r *Result) FQDN() string {
	parts := []string{}
	if r.Subdomain != "" {
		parts = append(parts, r.Subdomain)
	}
	if r.Domain != "" {
		parts = append(parts, r.Domain)
	}
	if r.TLD != "" {
		parts = append(parts, r.TLD)
	}
	return strings.Join(parts, ".")
}