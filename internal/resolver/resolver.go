package resolver

import (
	"errors"

	"dnstom/internal/dnswire"
)

type Resolver struct {
	server string // "host:port"
}

func New(server string) *Resolver {
	return &Resolver{server: server}
}

// Lookup is a generic entry point. You can later add helpers like LookupA, LookupAAAA, etc.
func (r *Resolver) Lookup(name, qtype string) (*dnswire.Message, error) {
	// For now, this is a stub so your project compiles.
	// Next steps:
	//   - build a dnswire.EncodeQuery(name, qtype)
	//   - send via net.DialUDP / net.Dial("udp", ...)
	//   - read response and decode into dnswire.Message
	return nil, errors.New("resolver.Lookup not implemented yet")
}
