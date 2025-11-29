package dnswire

import (
	"encoding/binary"
	"fmt"
	"io"
)

// PrintHeaderDiagram prints an RFC-style DNS header layout and the
// concrete values for the given Header.
func PrintHeaderDiagram(h Header, w io.Writer) error {
	// Encode the header so we can show raw values too.
	raw, err := EncodeHeader(h)
	if err != nil {
		return fmt.Errorf("encode header for diagram: %w", err)
	}

	if len(raw) != 12 {
		return fmt.Errorf("expected 12-byte DNS header, got %d bytes", len(raw))
	}

	id := binary.BigEndian.Uint16(raw[0:2])
	flags := binary.BigEndian.Uint16(raw[2:4])
	qd := binary.BigEndian.Uint16(raw[4:6])
	an := binary.BigEndian.Uint16(raw[6:8])
	ns := binary.BigEndian.Uint16(raw[8:10])
	ar := binary.BigEndian.Uint16(raw[10:12])

	// Decode flags field for pretty printing.
	qr := (flags>>15)&0x1 == 1
	opcode := (flags >> 11) & 0xF
	aa := (flags>>10)&0x1 == 1
	tc := (flags>>9)&0x1 == 1
	rd := (flags>>8)&0x1 == 1
	ra := (flags>>7)&0x1 == 1
	z := (flags >> 4) & 0x7
	rcode := flags & 0xF

	// RFC-style bit layout.
	if _, err := fmt.Fprintln(w, `
                                    1  1  1  1  1  1
      0  1  2  3  4  5  6  7  8  9  0  1  2  3  4  5
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |                      ID                       |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |QR|   Opcode  |AA|TC|RD|RA|   Z    |   RCODE   |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |                    QDCOUNT                    |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |                    ANCOUNT                    |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |                    NSCOUNT                    |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
     |                    ARCOUNT                    |
     +--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+
`); err != nil {
		return err
	}

	// Concrete values.
	if _, err := fmt.Fprintf(w, "ID:      0x%04X (%d)\n", id, id); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "FLAGS:   0x%04X\n", flags); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w,
		"  QR=%d  Opcode=%d  AA=%d  TC=%d  RD=%d  RA=%d  Z=%d  RCODE=%d\n",
		boolToBit(qr), opcode, boolToBit(aa), boolToBit(tc),
		boolToBit(rd), boolToBit(ra), z, rcode,
	); err != nil {
		return err
	}

	if _, err := fmt.Fprintf(w, "QDCOUNT: %d\nANCOUNT: %d\nNSCOUNT: %d\nARCOUNT: %d\n",
		qd, an, ns, ar); err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, "\nRaw header bytes (hex):"); err != nil {
		return err
	}
	if err := PrintRawHeaderBytes(raw, w); err != nil {
		return err
	}

	return nil
}

func boolToBit(b bool) int {
	if b {
		return 1
	}
	return 0
}

// PrintRawHeaderBytes prints the 12-byte header in a friendly hex format.
func PrintRawHeaderBytes(raw []byte, w io.Writer) error {
	if len(raw) != 12 {
		return fmt.Errorf("expected 12-byte DNS header, got %d bytes", len(raw))
	}

	_, err := fmt.Fprintf(w,
		"  %02X %02X  %02X %02X  %02X %02X  %02X %02X  %02X %02X  %02X %02X\n",
		raw[0], raw[1], raw[2], raw[3], raw[4], raw[5],
		raw[6], raw[7], raw[8], raw[9], raw[10], raw[11],
	)
	return err
}
