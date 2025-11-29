package dnswire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
)

const (
	// DNS record types
	TypeA uint16 = 1

	// DNS class
	ClassIN uint16 = 1
)

type Header struct {
	ID      uint16
	QR      bool
	Opcode  uint8
	AA      bool
	TC      bool
	RD      bool
	RA      bool
	Z       uint8 // usually 0
	Rcode   uint8
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type Question struct {
	Name  string
	Type  uint16
	Class uint16
}

type ResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte // you'll later decode into typed RRs using internal/rr
}

type Message struct {
	Header     Header
	Questions  []Question
	Answers    []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

// A Record:

// EncodeHeader encodes the DNS header into a 12-byte buffer.
func EncodeHeader(h Header) ([]byte, error) {
	// DNS header is always 12 bytes.
	var buf [12]byte

	// Use network byte order (big-endian).
	binary.BigEndian.PutUint16(buf[0:2], h.ID)

	// Flags: QR (1) | Opcode (4) | AA (1) | TC (1) | RD (1)
	//        RA (1) | Z (3)      | Rcode (4)
	var flags uint16 = 0

	if h.QR {
		flags |= 1 << 15
	}
	flags |= uint16(h.Opcode&0xF) << 11
	if h.AA {
		flags |= 1 << 10
	}
	if h.TC {
		flags |= 1 << 9
	}
	if h.RD {
		flags |= 1 << 8
	}
	if h.RA {
		flags |= 1 << 7
	}
	flags |= uint16(h.Z&0x7) << 4
	flags |= uint16(h.Rcode & 0xF)

	binary.BigEndian.PutUint16(buf[2:4], flags)
	binary.BigEndian.PutUint16(buf[4:6], h.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCount)

	return buf[:], nil
}

// EncodeQuery builds a DNS query message (header + question section).
// For this step, only the header encoding is implemented; the rest is TODO for you.
func EncodeQuery(name string, qtype uint16) ([]byte, error) {
	var h Header

	// ID: random 16-bit value
	h.ID = uint16(rand.Intn(65535))

	// This is a query
	h.QR = false

	// Standard query
	h.Opcode = 0

	// Recursion desired
	h.RD = true

	// One question
	h.QDCount = 1

	headerBytes, err := EncodeHeader(h)
	if err != nil {
		return nil, fmt.Errorf("encode header: %w", err)
	}

	// TODO: encode Question section:
	// - QNAME: split name into labels, prefix each with length byte, end with 0x00
	// - QTYPE: 2 bytes (e.g. TypeA)
	// - QCLASS: 2 bytes (ClassIN)

	var buf bytes.Buffer

	// Write header
	if _, err := buf.Write(headerBytes); err != nil {
		return nil, fmt.Errorf("write header: %w", err)
	}

	// TODO: encode QNAME based on `name` (e.g. "www.example.com")
	// Something like:
	//   labels := strings.Split(name, ".")
	//   for each label:
	//       write length byte
	//       write label bytes
	//   write terminating 0x00
	//
	// Then write QTYPE and QCLASS in big-endian.

	// Placeholder so function compiles; you'll replace this with real encoding.
	_ = qtype // to avoid unused warning for now

	return buf.Bytes(), nil
}

// PrettyPrint is a placeholder for now; you'll make this nicer later.
func PrettyPrint(m *Message, w io.Writer) error {
	_, err := fmt.Fprintf(w, ";; dnstom raw message\n;; ID: %d, QD: %d, AN: %d, NS: %d, AR: %d\n",
		m.Header.ID, m.Header.QDCount, m.Header.ANCount, m.Header.NSCount, m.Header.ARCount)
	if err != nil {
		return err
	}

	for _, q := range m.Questions {
		if _, err := fmt.Fprintf(w, ";; QUESTION: %s type=%d class=%d\n", q.Name, q.Type, q.Class); err != nil {
			return err
		}
	}

	for _, rr := range m.Answers {
		if _, err := fmt.Fprintf(w, "ANSWER: %s TTL=%d type=%d class=%d\n",
			rr.Name, rr.TTL, rr.Type, rr.Class); err != nil {
			return err
		}
	}

	return nil
}
