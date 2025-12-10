package resolver

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"time"

	"dnstom/internal/dnswire"
)

type Resolver struct {
	server  string // e.g. "1.1.1.1:53"
	diagram bool
}

func New(server string, diagram bool) *Resolver {
	return &Resolver{server: server}
}

// Called from main():
//	r := resolver.New(*server) -> create this.

//	ips, err := r.LookupA(name) -> lookup A Record.

// Test data passed in name =  example.com

// LookupA will eventually perform a real DNS query to r.server and return
// IPv4 addresses (A records) for the given name.
// Right now, EncodeQuery is called; UDP send/receive and response parsing
// are left as detailed TODOs for you.

func (r *Resolver) LookupA(name string) ([]net.IP, error) {
	// Step 1: Build a raw DNS query packet (header + question).
	query, err := dnswire.EncodeQuery(name, dnswire.TypeA)
	if err != nil {
		return nil, fmt.Errorf("build DNS query: %w", err)
	}

	// Query message output (for testing)
	fmt.Println("Hex dumper")
	d := hex.Dumper(os.Stdout)
	d.Write(query)
	d.Close()

	// Send query to name server

	nameserver := "8.8.8.8:53" // Replace with variable for parameter injection

	connection, err := net.Dial("udp", nameserver)
	if err != nil {
		panic(err)
	}

	defer connection.Close()

	// Timeout in case we don't get a response (prevent hanging)
	connection.SetDeadline((time.Now().Add(2 * time.Second)))

	// Send query packet
	_, err = connection.Write(query)
	if err != nil {
		panic(err)
	}

	// Buffer for the response - 512 is the max for DNS over UDP
	response := make([]byte, 512)

	n, err := connection.Read(response)
	if err != nil {
		panic(err)
	}

	// err, answer := dnswire.DecodeMessage(response)

	// To test

	fmt.Printf("Bytes received: %d   -    ", n)
	// fmt.Printf("%x\n", answer)

	// For now, just return an explicit “not implemented” error so you know
	// this code path is still a work in progress when it’s hit.
	return nil, fmt.Errorf("LookupA: DNS over UDP + response parsing not implemented yet")
}

// Optionally, later you might add a timeout configuration or helper:
//
// func (r *Resolver) timeout() time.Duration {
//     return 2 * time.Second
// }
