package dnswire

import (
	"encoding/binary"
	"fmt"
)

func decodeHeader(msg []byte) (header Header) {

	id := binary.BigEndian.Uint16(msg[0:2])
	flags := binary.BigEndian.Uint16(msg[2:4])

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

		QDCount: binary.BigEndian.Uint16(msg[4:6]),
		ANCount: binary.BigEndian.Uint16(msg[6:8]),
		NSCount: binary.BigEndian.Uint16(msg[8:10]),
		ARCount: binary.BigEndian.Uint16(msg[10:12]),
	}

	return hdr
}

// TODO: functionize this function
func decodeQuestion(questionBytes []byte) (question Question) {

	// Block to decode
	/* 		"7466"+"0100"+"0001"+"0000"+"0000"+"0000"+  // HEADER
	"03777777"+"057961686f6f"+"03636f6d"+"00"+ // QNAME
	"0001"+"0001", */ //QTYPE + QCLASS

	// binary.BigEndian.Uint16(msg[10:12]),
	fmt.Println("binary.BigEndian.Uint16(msg[10:12])")
	fmt.Printf("questionBytes[13]: %02X\n", questionBytes[13])
	fmt.Printf("questionBytes[14]: %02X\n", questionBytes[14])
	fmt.Printf("questionBytes[15]: %02X\n", questionBytes[15])
	fmt.Printf("questionBytes[16]: %02X\n", questionBytes[16])
	fmt.Printf("questionBytes[17]: %02X\n", questionBytes[17])
	fmt.Printf("questionBytes[18]: %02X\n", questionBytes[18])
	fmt.Printf("questionBytes[19]: %02X\n", questionBytes[19])

	fmt.Printf("questionBytes[20]: %02X\n", questionBytes[20])
	fmt.Printf("questionBytes[21]: %02X\n", questionBytes[21])

	fmt.Printf("questionBytes[22]: %02X\n", questionBytes[22])

	fmt.Printf("questionBytes[23]: %02X\n", questionBytes[23])

	q := Question{
		Name:  "name",
		Type:  TypeA,   //need to overwrite this with actual value
		Class: ClassIN, //need to overwrite this with actual value
	}
	return q
}

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
