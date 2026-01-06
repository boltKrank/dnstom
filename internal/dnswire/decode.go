package dnswire

import (
	"encoding/binary"
	"fmt"
	"strings"
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
func decodeQuestion(questionBytes []byte) (Question, int) {

	// Block to decode
	/* 		"7466"+"0100"+"0001"+"0000"+"0000"+"0000"+  // HEADER
	"03777777"+"057961686f6f"+"03636f6d"+"00"+ // QNAME
	"0001"+"0001", */ //QTYPE + QCLASS

	//Take off header
	offset := 12

	questionName, offset := decodeName(questionBytes, offset)

	fmt.Printf("\n Current offset after decodeName() is: %d \n", offset)

	q := Question{
		Name:  questionName,
		Type:  TypeA,   //need to overwrite this with actual value
		Class: ClassIN, //need to overwrite this with actual value
	}

	// Type (2 bytes) + Class (2 bytes)
	offset += 4 // 4 bytes total

	return q, offset
}

// TODO: Fille out this function
func decodeName(encodedName []byte, offset int) (string, int) {

	printByteArrayAsHex(encodedName)

	printByteArrayAsASCII(encodedName)

	// Decoding
	// 03 77 77 77 05 79 61 68 6f 6f 03 63 6f 6d 00
	// Check hex read those bytes and next hex is length again (unless 0x00)
	var labels []string

	fmt.Printf("\nTest - labels = %s", labels)

	for i := offset; i < len(encodedName) && encodedName[i] != 0; {
		l := int(encodedName[i])
		i++
		labels = append(labels, string(encodedName[i:i+l])) // Error thrown here dnstom/internal/dnswire.decodeName({0x140000c2160, 0x1f, 0x1f}) dnstom/internal/dnswire/decode.go:73 +0x244
		i += l
		offset++
	}

	decodedName := strings.Join(labels, ".")

	fmt.Println(decodedName)

	return decodedName, offset

}

func DecodeMessage(encodedMessage []byte) (Message, error) {

	// NOTE: Since the whole message is decoded here, we don't need to worry about the offset.
	// The offset's relevance finishes at the end of the message.

	/* 	The goal will be to return a struct that looks like this:
	   	type Message struct {
	   		Header     Header
	   		Questions  []Question
	   		Answers    []ResourceRecord
	   		Authority  []ResourceRecord
	   		Additional []ResourceRecord
	   	} */

	m := Message{
		Header:     decodeHeader(encodedMessage),
		Questions:  decodeQuestion(encodedMessage),
		Authority:  decodeAuthority(encodedMessage),
		Additional: decodeAdditional(encodedMessage),
	}

	return m
}

func decodeAdditional(encodedMessage []byte) ResourceRecord {
	panic("unimplemented")
}

func decodeAuthority(encodedMessage []byte) ResourceRecord {
	panic("unimplemented")
}

/*
	q, off2, err := decodeQuestion(packet []byte, off)
	if err != nil {
		t.Fatalf("decodeQuestion error: %v", err)
	} */

// Data received from www.example.com: Bytes received: 12   -    6caa81010000000000000000
