package dnswire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
)

const (
	// DNS record types
	TypeA uint16 = 1

	// DNS class
	ClassIN uint16 = 1
)

// Header, Question, ResourceRecord, Message, PrettyPrint stay as you already have them.

// EncodeHeader encodes the DNS header into a 12-byte buffer.
func EncodeHeader(h Header) ([]byte, error) {
	// DNS header is always 12 bytes.
	var buf [12]byte

	// Use network byte order (big-endian).
	binary.BigEndian.PutUint16(buf[0:2], h.ID)

	// Flags: QR (1) | Opcode (4) | AA (1) | TC (1) | RD (1)
	//        RA (1) | Z (3)      | RCODE (4)
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
// For this step, only the header encoding is implemented; the rest is
// left as TODOs with detailed guidance for you to fill in.
func EncodeQuery(name string, qtype uint16) ([]byte, error) {
	var h Header

	// ID: random 16-bit value (used to match response with query).
	h.ID = uint16(rand.Intn(65535))

	// This is a query (not a response).
	h.QR = false

	// Standard query (Opcode 0).
	h.Opcode = 0

	// Recursion desired: ask the upstream resolver to recurse for us.
	h.RD = true

	// One question in this message.
	h.QDCount = 1

	headerBytes, err := EncodeHeader(h)
	if err != nil {
		return nil, fmt.Errorf("encode header: %w", err)
	}

	var buf bytes.Buffer

	// Write the 12-byte header first.
	if _, err := buf.Write(headerBytes); err != nil {
		return nil, fmt.Errorf("write header: %w", err)
	}

	// TODO: Encode QNAME (the domain name in question)
	// Steps:
	//   1. Split the domain name into labels, e.g. "www.example.com" →
	//        []string{"www", "example", "com"}.
	//   2. For each label:
	//        - Check that len(label) <= 63 (DNS label length limit).
	//        - Write 1 byte: the length of the label.
	//        - Write N bytes: the ASCII bytes of the label.
	//   3. After all labels, write a zero length byte: 0x00 (root label terminator).
	// Notes:
	//   - Do NOT write the dots ('.') themselves, only length-prefixed labels.
	//   - QNAME always ends with 0x00, even for single-label names.
	//   - See RFC 1035 §3.1 “Name space definitions” for name encoding.
	//
	// Example encoding:
	//   "www.example.com" →
	//      03 'w' 'w' 'w' 07 'e' 'x' 'a' 'm' 'p' 'l' 'e' 03 'c' 'o' 'm' 00

	// TODO: Write QTYPE and QCLASS immediately after QNAME
	// Steps:
	//   1. QTYPE is a 16-bit unsigned integer in network byte order.
	//      - For A records, use TypeA (1).
	//   2. QCLASS is a 16-bit unsigned integer in network byte order.
	//      - For Internet, use ClassIN (1).
	//   3. Use binary.BigEndian.PutUint16 into a temporary [2]byte and write it
	//      into the buffer for both QTYPE and QCLASS.
	//
	// Resulting QUESTION layout:
	//   QNAME (variable length)
	//   QTYPE (2 bytes)
	//   QCLASS (2 bytes)
	//
	// Once done, buf.Bytes() will represent a complete DNS query message.

	_ = name  // to avoid unused warnings until you implement QNAME
	_ = qtype // to avoid unused warnings until you write QTYPE

	return buf.Bytes(), nil
}
