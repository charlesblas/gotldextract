# Go TLD Extract

A Go library and command-line tool for extracting the top-level domain (TLD), domain, and subdomain from URLs and domain names. This is a Go implementation inspired by the Python `tldextract` library.

## Features

- Extract TLD, domain, and subdomain from any URL or domain name
- Support for complex TLDs (e.g., `.co.uk`, `.com.br`)
- Handles URLs with protocols, ports, and paths
- Uses the Mozilla Public Suffix List for accurate TLD detection
- Command-line tool with support for stdin/stdout piping
- JSON output format option
- Automatic public suffix list updates

## Installation

### As a library

```bash
go get github.com/charlesblas/gotldextract
```

### As a command-line tool

```bash
go install github.com/charlesblas/gotldextract/cmd/tldextract@latest
```

## Usage

### Library Usage

```go
package main

import (
    "fmt"
    "github.com/charlesblas/gotldextract"
)

func main() {
    // Extract from a domain
    result, err := gotldextract.Extract("subdomain.example.com")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Subdomain:", result.Subdomain)  // "subdomain"
    fmt.Println("Domain:", result.Domain)        // "example"
    fmt.Println("TLD:", result.TLD)              // "com"
    fmt.Println("FQDN:", result.FQDN())          // "subdomain.example.com"
    
    // Extract from a URL
    result, _ = gotldextract.ExtractFromURL("https://www.example.co.uk/path")
    fmt.Println("Domain:", result.Domain)        // "example"
    fmt.Println("TLD:", result.TLD)              // "co.uk"
}
```

### Command-Line Usage

```bash
# Extract from command line arguments
tldextract example.com www.example.co.uk

# Extract from stdin
echo "www.example.com" | tldextract

# Extract from a file
tldextract -file domains.txt

# Output as JSON
tldextract -json example.com

# Update the public suffix list
tldextract -update

# Process a list of domains
cat domains.txt | tldextract > extracted.txt
```

### Command-Line Options

- `-update`: Update the public suffix list
- `-json`: Output results as JSON
- `-file`: Read from file instead of stdin
- `-help`: Show help message

## Examples

### Basic extraction
```bash
$ tldextract www.example.com
example.com
```

### JSON output
```bash
$ tldextract -json https://api.example.co.uk
{"input":"https://api.example.co.uk","subdomain":"api","domain":"example","tld":"co.uk","fqdn":"api.example.co.uk"}
```

### Batch processing
```bash
$ cat domains.txt
www.google.com
blog.example.co.uk
api.github.com

$ cat domains.txt | tldextract
google.com
example.co.uk
github.com
```

## API Reference

### `Extract(domain string) (*Result, error)`
Extracts the TLD, domain, and subdomain from a given input string.

### `ExtractFromURL(url string) (*Result, error)`
Alias for Extract that emphasizes it can handle full URLs.

### `Update() error`
Updates the local copy of the Mozilla Public Suffix List.

### `Result` struct
```go
type Result struct {
    Subdomain string  // The subdomain portion (e.g., "www", "api.v2")
    Domain    string  // The domain name (e.g., "example")
    TLD       string  // The top-level domain (e.g., "com", "co.uk")
}
```

Methods:
- `String() string`: Returns the registered domain (domain + TLD)
- `FQDN() string`: Returns the fully qualified domain name

## Performance

The library is optimized for high-performance batch processing:
- Uses buffered I/O for efficient file/stdin reading
- Minimal memory allocations
- Efficient string processing

## License

This project is open source and available under the MIT License.