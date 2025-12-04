package dnswire

import (
	"testing"
)

// TestPrintDNSDiagram creates a fully-populated DNSMessage and prints it.
// It does NOT validate correctness — it's for visual debug + development.
func TestPrintDNSDiagram(t *testing.T) {

	msg := DNSMessage{
		Header: DNSHeader{
			ID:      0x1234,
			Flags:   0x8180, // Standard "no error" DNS response
			QDCount: 1,
			ANCount: 1,
			NSCount: 0,
			ARCount: 0,
		},

		Questions: []DNSQuestion{
			{
				Name:   "www.example.com.",
				QType:  1, // A
				QClass: 1, // IN
			},
		},

		Answers: []DNSResourceRecord{
			{
				Name:     "www.example.com.",
				Type:     1,                        // A
				Class:    1,                        // IN
				TTL:      3600,                     // 1 hour
				RDLength: 4,                        // IPv4
				RData:    []byte{93, 184, 216, 34}, // 93.184.216.34
			},
		},

		Authorities: []DNSResourceRecord{}, // none
		Additionals: []DNSResourceRecord{}, // none
	}

	// This calls your RFC-style diagram printer.
	PrintDNSMessageDiagram(&msg)
}
