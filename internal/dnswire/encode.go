package dnswire

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"strings"
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

func EncodeQuestion(q Question) ([]byte, error) {

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

	/* func main() {
		qname := "www.google.com"

		labels := strings.Split(qname, ".")

		byte_size := 0

		for index, value := range labels {

			byte_size = byte_size + len(labels[index])
			fmt.Printf("Index: %d, Value: %s, Length: %d\n", index, value, len(labels[index]))
		}

		fmt.Printf("QNAME Length: %d", byte_size)

	} */

	labels := strings.Split(q.Name, ".")

	var qname []byte

	for _, label := range labels {

		qname = append(qname, byte(len(label))) //The byte length of the label coming up

		qname = append(qname, []byte(label)...) // The label itself.

	}

	//Finish off the QNAME with 0x00:
	qname = append(qname, 0x00)

	//Now we have the question, we just add the QTYPE and QCLASS
	question := qname

	// QTYPE
	question = append(question, byte(q.Type))

	// QCLASS
	question = append(question, byte(q.Class))

	return question, nil
}

// Header and question sections built here
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

	var q Question

	questionBytes, err := EncodeQuestion(q)
	if err != nil {
		return nil, fmt.Errorf("encode question: %w", err)
	}

	// Write the question section
	if _, err := buf.Write(questionBytes); err != nil {
		return nil, fmt.Errorf("write question: %w", err)
	}

	return buf.Bytes(), nil
}
