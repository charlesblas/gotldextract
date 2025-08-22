package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/charlesblas/gotldextract"
)

var (
	updateFlag = flag.Bool("update", false, "Update the public suffix list")
	jsonFlag   = flag.Bool("json", false, "Output results as JSON")
	fileFlag   = flag.String("file", "", "Input file (default: stdin)")
	helpFlag   = flag.Bool("help", false, "Show help")
)

type JSONOutput struct {
	Input     string `json:"input"`
	Subdomain string `json:"subdomain"`
	Domain    string `json:"domain"`
	TLD       string `json:"tld"`
	FQDN      string `json:"fqdn"`
}

func main() {
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}

	if *updateFlag {
		fmt.Println("Note: The public suffix list is embedded in the library.")
		fmt.Println("To update it, please update the golang.org/x/net/publicsuffix package:")
		fmt.Println("  go get -u golang.org/x/net/publicsuffix")
		return
	}

	// If there are command-line arguments (domains), process them
	if flag.NArg() > 0 {
		for _, domain := range flag.Args() {
			processDomain(domain)
		}
		return
	}

	// Otherwise, read from file or stdin
	var reader io.Reader
	if *fileFlag != "" {
		file, err := os.Open(*fileFlag)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	writer := bufio.NewWriterSize(os.Stdout, 64*1024)
	defer writer.Flush()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		processDomainWithWriter(line, writer)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading input: %v", err)
	}
}

func processDomain(domain string) {
	result, err := gotldextract.Extract(domain)
	if err != nil {
		log.Printf("Error processing %s: %v", domain, err)
		return
	}

	if *jsonFlag {
		output := JSONOutput{
			Input:     domain,
			Subdomain: result.Subdomain,
			Domain:    result.Domain,
			TLD:       result.TLD,
			FQDN:      result.FQDN(),
		}
		jsonBytes, _ := json.Marshal(output)
		fmt.Println(string(jsonBytes))
	} else {
		fmt.Printf("Input: %s\n", domain)
		fmt.Printf("  Subdomain: %s\n", result.Subdomain)
		fmt.Printf("  Domain: %s\n", result.Domain)
		fmt.Printf("  TLD: %s\n", result.TLD)
		fmt.Printf("  FQDN: %s\n", result.FQDN())
		fmt.Println()
	}
}

func processDomainWithWriter(domain string, writer *bufio.Writer) {
	result, err := gotldextract.Extract(domain)
	if err != nil {
		log.Printf("Error processing %s: %v", domain, err)
		return
	}

	if *jsonFlag {
		output := JSONOutput{
			Input:     domain,
			Subdomain: result.Subdomain,
			Domain:    result.Domain,
			TLD:       result.TLD,
			FQDN:      result.FQDN(),
		}
		jsonBytes, _ := json.Marshal(output)
		writer.WriteString(string(jsonBytes))
		writer.WriteByte('\n')
	} else {
		// For pipe mode, just output the extracted domain.tld
		writer.WriteString(result.String())
		writer.WriteByte('\n')
	}
}

func printHelp() {
	fmt.Println("tldextract - Extract TLD, domain, and subdomain from URLs and domain names")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  tldextract [options] [domains...]")
	fmt.Println("  command | tldextract [options]")
	fmt.Println("  tldextract [options] < domains.txt")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -update    Update the public suffix list")
	fmt.Println("  -json      Output results as JSON")
	fmt.Println("  -file      Read from file instead of stdin")
	fmt.Println("  -help      Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  tldextract example.com")
	fmt.Println("  tldextract -json https://www.example.co.uk")
	fmt.Println("  echo 'subdomain.example.com' | tldextract")
	fmt.Println("  tldextract -file domains.txt")
}