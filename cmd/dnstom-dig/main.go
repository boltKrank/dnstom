package main

import (
	"flag" //CLI flags
	"fmt"  // Strings
	"log"
	"os" // Seems like unless it's OS/2 it's probably covered :P

	"dnstom/internal/resolver"
)

func main() {
	// Print out RFC-style diagrams ?
	diagram := flag.Bool("diagram", true, "Print RFC-Style diagrams for network packets")

	// We accept a --server flag for future use, even though step 1 ignores it.
	server := flag.String("server", "system", "DNS server to query")
	flag.Parse()

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "dnstom-dig - a toy DNS resolver\n")
		fmt.Fprintf(os.Stderr, "Usage: dnstom-dig [options] <name>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	if flag.NArg() < 1 {
		flag.Usage()
		return
	}

	name := flag.Arg(0)

	r := resolver.New(*server, *diagram)

	// We're going to be looking up an A record here.
	// IPv6 might be implemented later, but I don't know how long palliative care will continue
	// for that dying beast. (Although it keeps SecOps in business)
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
