package resolver

import (
	"fmt"
	"net"

	"dnstom/internal/dnswire"
)

type Resolver struct {
	server string // e.g. "1.1.1.1:53"
}

func New(server string) *Resolver {
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

	// TODO: Step 2 — Send the query via UDP to r.server and receive a response
	//
	// Suggested steps:
	//   1. Dial a UDP connection to the configured DNS server:
	//        conn, err := net.Dial("udp", r.server)
	//      - r.server should be something like "1.1.1.1:53".
	//      - Handle the error if dialing fails.
	//
	//   2. Optionally set a deadline so the read does not block forever:
	//        deadline := time.Now().Add(2 * time.Second)
	//        err = conn.SetDeadline(deadline)
	//      - Choose a timeout you consider reasonable for your lab.
	//
	//   3. Write the query bytes to the connection:
	//        _, err = conn.Write(query)
	//      - Ensure you handle short writes or errors.
	//
	//   4. Allocate a buffer for the response:
	//        buf := make([]byte, 512)
	//      - 512 bytes is the classic DNS-over-UDP limit without EDNS0.
	//      - Later you can support larger messages via EDNS.
	//
	//   5. Read from the connection:
	//        n, err := conn.Read(buf)
	//      - Handle read errors (timeout, connection issues, etc.).
	//      - Use buf[:n] as the actual response bytes.
	//
	//   6. Close the connection when done:
	//        conn.Close()
	//
	// Notes:
	//   - For now, you can ignore retransmits, truncation (TC flag), and
	//     other robustness concerns; those can come later.

	_ = query // remove this once you actually send it

	// TODO: Step 3 — Decode the response bytes into a dnswire.Message
	//
	// Steps:
	//   1. Implement a function like:
	//        func DecodeMessage(raw []byte) (*dnswire.Message, error)
	//      in the dnswire package.
	//   2. Inside DecodeMessage:
	//        - Parse the first 12 bytes into a Header.
	//        - Check that the ID matches the query ID (to guard against
	//          mismatched or out-of-order responses in future).
	//        - Parse questions, answers, authority, and additional sections.
	//        - Store results in a dnswire.Message struct.
	//   3. Call that function here and handle any errors.
	//
	// For the first version, you can:
	//   - Only fully decode the header and answer section.
	//   - Focus just on A records in the answer section.
	//   - Leave DNS compression and more complex RR types for later steps.

	// TODO: Step 4 — Extract A records from the Answer section into []net.IP
	//
	// Steps (once you have a decoded dnswire.Message):
	//   1. Iterate over msg.Answers.
	//   2. For each ResourceRecord:
	//        - Check if rr.Type == dnswire.TypeA and rr.Class == dnswire.ClassIN.
	//        - Ensure rr.RData has length 4 (IPv4 address).
	//        - Convert the 4 bytes into a net.IP:
	//             ip := net.IP(rr.RData)
	//        - Append ip to a slice of net.IP.
	//   3. Return the slice of net.IP and a nil error if everything succeeds.
	//
	// Later improvements:
	//   - Handle cases where the answer is a CNAME pointing to another name.
	//   - Follow up queries for that CNAME target (re-query or use ADDED records).
	//   - Handle multiple A records (round-robin, etc.).

	// For now, just return an explicit “not implemented” error so you know
	// this code path is still a work in progress when it’s hit.
	return nil, fmt.Errorf("LookupA: DNS over UDP + response parsing not implemented yet")
}

// Optionally, later you might add a timeout configuration or helper:
//
// func (r *Resolver) timeout() time.Duration {
//     return 2 * time.Second
// }
