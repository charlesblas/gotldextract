package gotldextract

import (
	"testing"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		subdomain string
		domain    string
		tld       string
		fqdn      string
	}{
		{
			name:      "simple domain",
			input:     "example.com",
			subdomain: "",
			domain:    "example",
			tld:       "com",
			fqdn:      "example.com",
		},
		{
			name:      "subdomain",
			input:     "www.example.com",
			subdomain: "www",
			domain:    "example",
			tld:       "com",
			fqdn:      "www.example.com",
		},
		{
			name:      "multiple subdomains",
			input:     "a.b.c.example.com",
			subdomain: "a.b.c",
			domain:    "example",
			tld:       "com",
			fqdn:      "a.b.c.example.com",
		},
		{
			name:      "UK domain",
			input:     "example.co.uk",
			subdomain: "",
			domain:    "example",
			tld:       "co.uk",
			fqdn:      "example.co.uk",
		},
		{
			name:      "UK domain with subdomain",
			input:     "www.example.co.uk",
			subdomain: "www",
			domain:    "example",
			tld:       "co.uk",
			fqdn:      "www.example.co.uk",
		},
		{
			name:      "URL with protocol",
			input:     "https://www.example.com",
			subdomain: "www",
			domain:    "example",
			tld:       "com",
			fqdn:      "www.example.com",
		},
		{
			name:      "URL with path",
			input:     "https://www.example.com/path/to/page",
			subdomain: "www",
			domain:    "example",
			tld:       "com",
			fqdn:      "www.example.com",
		},
		{
			name:      "URL with port",
			input:     "https://www.example.com:8080",
			subdomain: "www",
			domain:    "example",
			tld:       "com",
			fqdn:      "www.example.com",
		},
		{
			name:      "Japanese domain",
			input:     "example.co.jp",
			subdomain: "",
			domain:    "example",
			tld:       "co.jp",
			fqdn:      "example.co.jp",
		},
		{
			name:      "Brazilian domain",
			input:     "example.com.br",
			subdomain: "",
			domain:    "example",
			tld:       "com.br",
			fqdn:      "example.com.br",
		},
		{
			name:      "Just TLD",
			input:     "com",
			subdomain: "",
			domain:    "",
			tld:       "com",
			fqdn:      "com",
		},
		{
			name:      "Mixed case",
			input:     "WWW.EXAMPLE.COM",
			subdomain: "www",
			domain:    "example",
			tld:       "com",
			fqdn:      "www.example.com",
		},
		{
			name:      "Trailing dot",
			input:     "example.com.",
			subdomain: "",
			domain:    "example",
			tld:       "com",
			fqdn:      "example.com",
		},
		{
			name:      "Complex subdomain",
			input:     "api.v2.staging.example.com",
			subdomain: "api.v2.staging",
			domain:    "example",
			tld:       "com",
			fqdn:      "api.v2.staging.example.com",
		},
		{
			name:      "GitHub Pages domain",
			input:     "username.github.io",
			subdomain: "",
			domain:    "username",
			tld:       "github.io",
			fqdn:      "username.github.io",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Extract(tt.input)
			if err != nil {
				t.Fatalf("Extract() error = %v", err)
			}

			if result.Subdomain != tt.subdomain {
				t.Errorf("Subdomain = %v, want %v", result.Subdomain, tt.subdomain)
			}
			if result.Domain != tt.domain {
				t.Errorf("Domain = %v, want %v", result.Domain, tt.domain)
			}
			if result.TLD != tt.tld {
				t.Errorf("TLD = %v, want %v", result.TLD, tt.tld)
			}
			if result.FQDN() != tt.fqdn {
				t.Errorf("FQDN() = %v, want %v", result.FQDN(), tt.fqdn)
			}
		})
	}
}

func TestExtractFromURL(t *testing.T) {
	// Test that ExtractFromURL works the same as Extract
	url := "https://www.example.com/path"
	result1, err1 := Extract(url)
	result2, err2 := ExtractFromURL(url)

	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v, %v", err1, err2)
	}

	if result1.Subdomain != result2.Subdomain ||
		result1.Domain != result2.Domain ||
		result1.TLD != result2.TLD {
		t.Errorf("ExtractFromURL produced different result than Extract")
	}
}

func TestResultString(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		want   string
	}{
		{
			name: "normal domain",
			result: Result{
				Subdomain: "www",
				Domain:    "example",
				TLD:       "com",
			},
			want: "example.com",
		},
		{
			name: "just TLD",
			result: Result{
				Subdomain: "",
				Domain:    "",
				TLD:       "com",
			},
			want: "com",
		},
		{
			name: "multi-part TLD",
			result: Result{
				Subdomain: "www",
				Domain:    "example",
				TLD:       "co.uk",
			},
			want: "example.co.uk",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.result.String(); got != tt.want {
				t.Errorf("Result.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCleanDomain(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", "example.com", "example.com"},
		{"with http", "http://example.com", "example.com"},
		{"with https", "https://example.com", "example.com"},
		{"with path", "example.com/path", "example.com"},
		{"with port", "example.com:8080", "example.com"},
		{"trailing dot", "example.com.", "example.com"},
		{"uppercase", "EXAMPLE.COM", "example.com"},
		{"complex", "HTTPS://EXAMPLE.COM:8080/PATH", "example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanDomain(tt.input); got != tt.want {
				t.Errorf("cleanDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}