package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BoltKrank/dnstom/internal/dnswire"
	"github.com/BoltKrank/dnstom/internal/resolver"
)

func main() {
	server := flag.String("server", "1.1.1.1:53", "DNS server to query (host:port)")
	qtype := flag.String("type", "A", "Query type (A, AAAA, MX, NS, TXT, etc.)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: dnstom-dig [options] <name>\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	name := flag.Arg(0)

	r := resolver.New(*server)

	msg, err := r.Lookup(name, *qtype)
	if err != nil {
		log.Fatalf("lookup error: %v", err)
	}

	if err := dnswire.PrettyPrint(msg, os.Stdout); err != nil {
		log.Fatalf("print error: %v", err)
	}
}
