package main

import (
	"flag" //CLI flags
	"fmt"  // Strings
	"log"
	"os"

	"dnstom/internal/resolver"
)

func main() {
	// We accept a --server flag for future use, even though step 1 ignores it.
	server := flag.String("server", "system", "DNS server to query (ignored in step 1; uses system resolver)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: dnstom-dig [options] <name>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	name := flag.Arg(0)

	r := resolver.New(*server)

	ips, err := r.LookupA(name)
	if err != nil {
		log.Fatalf("lookup error: %v", err)
	}

	if len(ips) == 0 {
		fmt.Printf("No IPv4 addresses found for %s\n", name)
		return
	}

	fmt.Printf("A records for %s:\n", name)
	for _, ip := range ips {
		fmt.Printf("  %s\n", ip.String())
	}
}
