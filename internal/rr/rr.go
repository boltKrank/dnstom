package rr

// Later: typed representations for A, AAAA, NS, MX, SOA, etc.

type A struct {
	Address string // you might use net.IP later
}

// You might eventually have:
// type SOA struct { ... }
// type MX struct { ... }
