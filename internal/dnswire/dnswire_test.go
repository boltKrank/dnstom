package dnswire

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// Expecting these constants are defined within the package:

/* const (
	TypeA   uint16 = 1
	ClassIN uint16 = 1
) */

// These should be declared in the message file

/* type Header struct {
	ID      uint16
	QR      bool
	Opcode  uint8
	AA      bool
	TC      bool
	RD      bool
	RA      bool
	Z       uint8
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
} */

// ---------- Helpers for tests ----------

func mustHexDecode(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("failed to hex decode %q: %v", s, err)
	}
	return b
}

// ---------- encodeName / decodeName ----------

func TestEncodeName_WWWYahooCom(t *testing.T) {
	name := "www.yahoo.com"

	got, err := encodeName(name)

	if err != nil {
		t.Fatalf("encodeName returned an error: %v", err)
	}

	// 03 'w' 'w' 'w' 05 'y' ... 'o' 'o' 03 'c' 'o' 'm' 00
	want := []byte{
		0x03, 'w', 'w', 'w',
		0x05, 'y', 'a', 'h', 'o', 'o',
		0x03, 'c', 'o', 'm',
		0x00,
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("encodeName(%q) = %x, want %x", name, got, want)
	}
}

func TestDecodeName_WWWYahooCom(t *testing.T) {
	encoded := []byte{
		0x03, 'w', 'w', 'w',
		0x05, 'y', 'a', 'h', 'o', 'o',
		0x03, 'c', 'o', 'm',
		0x00,
	}

	name, off, err := decodeName(encoded, 0)
	if err != nil {
		t.Fatalf("decodeName returned error: %v", err)
	}

	if name != "www.yahoo.com" {
		t.Fatalf("decodeName name = %q, want %q", name, "www.yahoo.com")
	}

	if off != len(encoded) {
		t.Fatalf("decodeName offset = %d, want %d", off, len(encoded))
	}
}

// ---------- encodeHeader / decodeHeader ----------

func TestEncodeHeader_BasicQuery(t *testing.T) {
	h := Header{
		ID:      0x7466,
		QR:      false,
		Opcode:  0,
		AA:      false,
		TC:      false,
		RD:      true,
		RA:      false,
		Z:       0,
		Rcode:   0,
		QDCount: 1,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	got, err := encodeHeader(h)

	if err != nil {
		t.Fatalf("encodeHeader returned an error: %v", err)
	}

	// Expected header bytes:
	// ID      = 0x7466
	// FLAGS   = 0x0100 (standard query, RD=1)
	// QDCOUNT = 0x0001
	// ANCOUNT = 0x0000
	// NSCOUNT = 0x0000
	// ARCOUNT = 0x0000
	want := mustHexDecode(t, "7466"+"0100"+"0001"+"0000"+"0000"+"0000")

	if !bytes.Equal(got, want) {
		t.Fatalf("encodeHeader() = %x, want %x", got, want)
	}
}

func TestDecodeHeader_BasicQuery(t *testing.T) {
	packet := mustHexDecode(t, "7466"+"0100"+"0001"+"0000"+"0000"+"0000")

	h, off, err := decodeHeader(packet)
	if err != nil {
		t.Fatalf("decodeHeader error: %v", err)
	}

	if off != len(packet) {
		t.Fatalf("decodeHeader offset = %d, want %d", off, len(packet))
	}

	if h.ID != 0x7466 {
		t.Errorf("ID = 0x%04x, want 0x7466", h.ID)
	}
	if h.QR != false {
		t.Errorf("QR = %v, want false", h.QR)
	}
	if h.Opcode != 0 {
		t.Errorf("Opcode = %d, want 0", h.Opcode)
	}
	if h.RD != true {
		t.Errorf("RD = %v, want true", h.RD)
	}
	if h.Rcode != 0 {
		t.Errorf("Rcode = %d, want 0", h.Rcode)
	}
	if h.QDCount != 1 {
		t.Errorf("QDCount = %d, want 1", h.QDCount)
	}
	if h.ANCount != 0 || h.NSCount != 0 || h.ARCount != 0 {
		t.Errorf("AN/NS/AR counts = %d/%d/%d, want 0/0/0", h.ANCount, h.NSCount, h.ARCount)
	}
}

// ---------- encodeQuestion / decodeQuestion + full packet ----------

func TestEncodeFullQuery_WWWYahooCom_A_IN(t *testing.T) {
	h := Header{
		ID:      0x7466,
		QR:      false,
		Opcode:  0,
		AA:      false,
		TC:      false,
		RD:      true,
		RA:      false,
		Z:       0,
		Rcode:   0,
		QDCount: 1,
		ANCount: 0,
		NSCount: 0,
		ARCount: 0,
	}

	q := Question{
		Name:  "www.yahoo.com",
		Type:  TypeA,
		Class: ClassIN,
	}

	var packet []byte
	packet = append(packet, encodeHeader(h)...)
	packet = append(packet, encodeQuestion(q)...)

	// Expected full query packet:
	// 74 66 01 00 00 01 00 00 00 00 00 00
	// 03 77 77 77 05 79 61 68 6f 6f 03 63 6f 6d 00
	// 00 01 00 01
	want := mustHexDecode(t,
		"7466"+"0100"+"0001"+"0000"+"0000"+"0000"+
			"03777777"+"057961686f6f"+"03636f6d"+"00"+
			"0001"+"0001",
	)

	if !bytes.Equal(packet, want) {
		t.Fatalf("full query packet = %x, want %x", packet, want)
	}
}

func TestDecodeFullQuery_WWWYahooCom_A_IN(t *testing.T) {
	// Same bytes as above test
	packet := mustHexDecode(t,
		"7466"+"0100"+"0001"+"0000"+"0000"+"0000"+
			"03777777"+"057961686f6f"+"03636f6d"+"00"+
			"0001"+"0001",
	)

	h, off, err := decodeHeader(packet)
	if err != nil {
		t.Fatalf("decodeHeader error: %v", err)
	}

	if h.ID != 0x7466 {
		t.Errorf("ID = 0x%04x, want 0x7466", h.ID)
	}
	if h.RD != true {
		t.Errorf("RD = %v, want true", h.RD)
	}
	if h.QDCount != 1 {
		t.Fatalf("QDCount = %d, want 1", h.QDCount)
	}

	q, off2, err := decodeQuestion(packet, off)
	if err != nil {
		t.Fatalf("decodeQuestion error: %v", err)
	}

	if q.Name != "www.yahoo.com" {
		t.Errorf("Question.Name = %q, want %q", q.Name, "www.yahoo.com")
	}
	if q.Type != TypeA {
		t.Errorf("Question.Type = %d, want %d", q.Type, TypeA)
	}
	if q.Class != ClassIN {
		t.Errorf("Question.Class = %d, want %d", q.Class, ClassIN)
	}

	if off2 != len(packet) {
		t.Errorf("final offset = %d, want %d", off2, len(packet))
	}
}
