package main

import (
	"fmt"
	"log"

	"github.com/charlesblas/gotldextract"
)

func main() {
	// Example domains to process
	domains := []string{
		"www.example.com",
		"https://blog.example.co.uk/post",
		"subdomain.example.github.io",
		"api.v2.staging.mycompany.com",
		"example.com.br",
		"mail.google.com",
	}

	fmt.Println("Go TLD Extract - Example Usage")
	fmt.Println("==============================")
	fmt.Println()

	for _, domain := range domains {
		fmt.Printf("Processing: %s\n", domain)
		
		result, err := gotldextract.Extract(domain)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("  Subdomain: %s\n", result.Subdomain)
		fmt.Printf("  Domain:    %s\n", result.Domain)
		fmt.Printf("  TLD:       %s\n", result.TLD)
		fmt.Printf("  Full:      %s\n", result.FQDN())
		fmt.Printf("  Registered Domain: %s\n", result.String())
		fmt.Println()
	}

	// Note about updates
	fmt.Println("Note: The public suffix list is embedded in the library.")
	fmt.Println("To get the latest TLD data, update the package with:")
	fmt.Println("  go get -u golang.org/x/net/publicsuffix")
}