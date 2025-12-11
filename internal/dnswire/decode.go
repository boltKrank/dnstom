package dnswire

import "encoding/binary"

func decodeHeader(msg []byte) Header {
	id := binary.BigEndian.Uint16(msg[0:2])
	flags := binary.BigEndian.Uint16(msg[2:4])

	hdr := Header{
		ID: id,

		QR:     (flags>>15)&0x1 == 1,
		Opcode: uint8((flags >> 11) & 0xF),
		AA:     (flags>>10)&0x1 == 1,
		TC:     (flags>>9)&0x1 == 1,
		RD:     (flags>>8)&0x1 == 1,
		RA:     (flags>>7)&0x1 == 1,

		Z:     uint8((flags >> 4) & 0x7),
		Rcode: uint8(flags & 0xF),

		QDCount: binary.BigEndian.Uint16(msg[4:6]),
		ANCount: binary.BigEndian.Uint16(msg[6:8]),
		NSCount: binary.BigEndian.Uint16(msg[8:10]),
		ARCount: binary.BigEndian.Uint16(msg[10:12]),
	}

	return hdr
}


// TODO: functionize this function
func decodeQuestion()

	q, off2, err := decodeQuestion(packet []byte, off)
	if err != nil {
		t.Fatalf("decodeQuestion error: %v", err)
	}


// Data received from www.example.com: Bytes received: 12   -    6caa81010000000000000000
