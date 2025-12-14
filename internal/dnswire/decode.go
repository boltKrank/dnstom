package dnswire

import (
	"encoding/binary"
	"fmt"
)

func decodeHeader(msg []byte) {
	if len(msg) < off+12 {
		return Header{}, off, fmt.Errorf("short DNS message: need %d bytes for header", off+12)
	}

	id := binary.BigEndian.Uint16(msg[off : off+2])
	flags := binary.BigEndian.Uint16(msg[off+2 : off+4])

	hdr := Header{
		ID: id,

		QR:     (flags>>15)&1 == 1,
		Opcode: uint8((flags >> 11) & 0xF),
		AA:     (flags>>10)&1 == 1,
		TC:     (flags>>9)&1 == 1,
		RD:     (flags>>8)&1 == 1,
		RA:     (flags>>7)&1 == 1,

		Z:     uint8((flags >> 4) & 0x7),
		Rcode: uint8(flags & 0xF),

		QDCount: binary.BigEndian.Uint16(msg[off+4 : off+6]),
		ANCount: binary.BigEndian.Uint16(msg[off+6 : off+8]),
		NSCount: binary.BigEndian.Uint16(msg[off+8 : off+10]),
		ARCount: binary.BigEndian.Uint16(msg[off+10 : off+12]),
	}

	return hdr, off + 12, nil
}

// TODO: functionize this function
func decodeQuestion(question []byte, off int) {}

// TODO: Fille out this function
func decodeName(encodedName []byte, offset int) (string, int, error) {

	return "not yet", 69, nil

}

/*
	q, off2, err := decodeQuestion(packet []byte, off)
	if err != nil {
		t.Fatalf("decodeQuestion error: %v", err)
	} */

// Data received from www.example.com: Bytes received: 12   -    6caa81010000000000000000
