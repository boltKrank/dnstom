package resolver

import (
	"net"
)

// Resolver represents a DNS resolver. For now, it just wraps net.LookupIP.
// Later you can choose to honour the server field if you want.
type Resolver struct {
	server string // e.g. "1.1.1.1:53" – unused in step 1
}

// New creates a new Resolver. The server argument is accepted for future use.
func New(server string) *Resolver {
	return &Resolver{server: server}
}

// LookupA resolves IPv4 addresses (A records) for the given hostname
// using the system resolver via net.LookupIP.
func (r *Resolver) LookupA(name string) ([]net.IP, error) {
	ips, err := net.LookupIP(name)
	if err != nil {
		return nil, err
	}

	var v4 []net.IP
	for _, ip := range ips {
		if ip.To4() != nil {
			v4 = append(v4, ip)
		}
	}
	return v4, nil
}
