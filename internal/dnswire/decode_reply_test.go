package dnswire

import (
	"encoding/hex"
	"net"
	"testing"
)

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(stripSpaces(s))
	if err != nil {
		t.Fatalf("bad hex fixture: %v", err)
	}
	return b
}

func stripSpaces(s string) string {
	out := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\n', '\r', '\t':
			continue
		default:
			out = append(out, s[i])
		}
	}
	return string(out)
}

func TestDecodeReply_ARecord_WithCompressionPointer(t *testing.T) {
	// Known-good DNS response (includes NAME compression pointer c00c in Answer NAME)
	// Response for: A www.northeastern.edu
	wire := mustHex(t, `
		db42 8180 0001 0001 0000 0000
		0377 7777 0c6e 6f72 7468 6561 7374 6572 6e03 6564 7500
		0001 0001
		c00c 0001 0001 0000 0258 0004 9b21 1144
	`)

	m, err := DecodeMessage(wire)
	if err != nil {
		t.Fatalf("DecodeMessage error: %v", err)
	}

	// ---- Header ----
	if m.Header.ID != 0xdb42 {
		t.Fatalf("ID: got 0x%04x want 0xdb42", m.Header.ID)
	}
	if !m.Header.QR {
		t.Fatalf("QR: got false want true (response)")
	}
	if m.Header.QDCount != 1 || m.Header.ANCount != 1 || m.Header.NSCount != 0 || m.Header.ARCount != 0 {
		t.Fatalf("counts: QD=%d AN=%d NS=%d AR=%d want 1,1,0,0",
			m.Header.QDCount, m.Header.ANCount, m.Header.NSCount, m.Header.ARCount)
	}

	// ---- Question ----
	if len(m.Questions) != 1 {
		t.Fatalf("Questions: got %d want 1", len(m.Questions))
	}
	q := m.Questions[0]
	if q.Name != "www.northeastern.edu." {
		t.Fatalf("QNAME: got %q want %q", q.Name, "www.northeastern.edu.")
	}
	if q.Type != 1 {
		t.Fatalf("QTYPE: got %d want 1 (A)", q.Type)
	}
	if q.Class != 1 {
		t.Fatalf("QCLASS: got %d want 1 (IN)", q.Class)
	}

	// ---- Answer ----
	if len(m.Answers) != 1 {
		t.Fatalf("Answers: got %d want 1", len(m.Answers))
	}
	a := m.Answers[0]
	// Name is a compression pointer to offset 0x000c, which is the start of QNAME
	if a.Name != "www.northeastern.edu." {
		t.Fatalf("ANAME: got %q want %q", a.Name, "www.northeastern.edu.")
	}
	if a.Type != 1 {
		t.Fatalf("ATYPE: got %d want 1 (A)", a.Type)
	}
	if a.Class != 1 {
		t.Fatalf("ACLASS: got %d want 1 (IN)", a.Class)
	}
	if a.TTL != 600 {
		t.Fatalf("TTL: got %d want 600", a.TTL)
	}
	if a.RDLength != 4 {
		t.Fatalf("RDLength: got %d want 4", a.RDLength)
	}
	if len(a.RData) != 4 {
		t.Fatalf("RData len: got %d want 4", len(a.RData))
	}

	gotIP := net.IP(a.RData)
	wantIP := net.IPv4(155, 33, 17, 68) // 9b 21 11 44
	if !gotIP.Equal(wantIP) {
		t.Fatalf("A RDATA: got %v want %v", gotIP, wantIP)
	}

	// Sanity: other sections empty
	if len(m.Authority) != 0 || len(m.Additional) != 0 {
		t.Fatalf("expected Authority/Additional empty, got NS=%d AR=%d", len(m.Authority), len(m.Additional))
	}
}

func TestDecodeReply_Truncated_ShouldError(t *testing.T) {
	wire := mustHex(t, `db42 8180 0001`)
	_, err := DecodeMessage(wire)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestDecodeReply_BadCompressionPointer_OutOfRange_ShouldError(t *testing.T) {
	// Header: QD=1, rest 0; QNAME starts with a pointer to 0x00ff (past end)
	wire := mustHex(t, `
		db42 0100 0001 0000 0000 0000
		c0ff 0001 0001
	`)
	_, err := DecodeMessage(wire)
	if err == nil {
		t.Fatalf("expected error due to invalid compression pointer, got nil")
	}
}
