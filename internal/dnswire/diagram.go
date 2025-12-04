// diagram.go
package dnswire // change to your package name

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// =============[ ANSI COLOR SETUP ]=============

// Toggle this if you want to disable colours globally.
var UseColor = true

const (
	cReset  = "\033[0m"
	cDim    = "\033[2m"
	cBold   = "\033[1m"
	cCyan   = "\033[36m"
	cWhite  = "\033[97m"
	cYellow = "\033[93m"
	cGreen  = "\033[92m"
)

func col(s, color string) string {
	if !UseColor {
		return s
	}
	return color + s + cReset
}

// =============[ TYPES ]=============

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

type DNSQuestion struct {
	Name   string // e.g. "www.example.com."
	QType  uint16
	QClass uint16
}

type DNSResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte
}

type DNSMessage struct {
	Header      DNSHeader
	Questions   []DNSQuestion
	Answers     []DNSResourceRecord
	Authorities []DNSResourceRecord
	Additionals []DNSResourceRecord
}

// Public entry point
func PrintDNSMessageDiagram(msg *DNSMessage) {
	printHeaderDiagram(&msg.Header)

	if len(msg.Questions) > 0 {
		fmt.Println()
		fmt.Println(col(";; QUESTION SECTION:", cCyan+cBold))
		for i, q := range msg.Questions {
			printQuestionDiagram(i, &q)
		}
	}

	if len(msg.Answers) > 0 {
		fmt.Println()
		fmt.Println(col(";; ANSWER SECTION:", cCyan+cBold))
		for i, rr := range msg.Answers {
			printRRDiagram(i, &rr)
		}
	}

	if len(msg.Authorities) > 0 {
		fmt.Println()
		fmt.Println(col(";; AUTHORITY SECTION:", cCyan+cBold))
		for i, rr := range msg.Authorities {
			printRRDiagram(i, &rr)
		}
	}

	if len(msg.Additionals) > 0 {
		fmt.Println()
		fmt.Println(col(";; ADDITIONAL SECTION:", cCyan+cBold))
		for i, rr := range msg.Additionals {
			printRRDiagram(i, &rr)
		}
	}
}

// =============[ HEADER ]=============

func printHeaderDiagram(h *DNSHeader) {
	flags := h.Flags

	qr := (flags >> 15) & 0x1
	opcode := (flags >> 11) & 0xF
	aa := (flags >> 10) & 0x1
	tc := (flags >> 9) & 0x1
	rd := (flags >> 8) & 0x1
	ra := (flags >> 7) & 0x1
	z := (flags >> 6) & 0x1
	ad := (flags >> 5) & 0x1
	cd := (flags >> 4) & 0x1
	rcode := flags & 0xF

	box := col("+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+", cDim)
	fmt.Println(col(";; HEADER", cCyan+cBold))
	fmt.Println(box)
	fmt.Println(col("|                      ID                       |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("|                   0x%04x                     |", h.ID), cYellow))
	fmt.Println(box)
	fmt.Println(col("|QR|   Opcode  |AA|TC|RD|RA| Z|AD|CD|   RCODE   |", cWhite))
	fmt.Printf("%s\n", col(
		fmt.Sprintf("| %d|   %4d   | %d| %d| %d| %d| %d| %d| %d|   %4d   |",
			qr, opcode, aa, tc, rd, ra, z, ad, cd, rcode),
		cYellow,
	))
	fmt.Println(box)

	fmt.Println(col("|                    QDCOUNT                    |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("|                   %10d                 |", h.QDCount), cYellow))
	fmt.Println(box)

	fmt.Println(col("|                    ANCOUNT                    |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("|                   %10d                 |", h.ANCount), cYellow))
	fmt.Println(box)

	fmt.Println(col("|                    NSCOUNT                    |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("|                   %10d                 |", h.NSCount), cYellow))
	fmt.Println(box)

	fmt.Println(col("|                    ARCOUNT                    |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("|                   %10d                 |", h.ARCount), cYellow))
	fmt.Println(box)
}

// =============[ QUESTION SECTION ]=============

func printQuestionDiagram(index int, q *DNSQuestion) {
	box := col("+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+", cDim)

	fmt.Println(box)
	fmt.Println(col(fmt.Sprintf(";; Question %d", index+1), cCyan+cBold))
	fmt.Println(col("|                     QNAME                     |", cWhite))
	fmt.Printf("%s\n", col("| "+padRight(q.Name, 45), cYellow))
	fmt.Println(box)
	fmt.Println(col("|                    QTYPE                      |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("| %-5d (%s)", q.QType, typeToString(q.QType)), cYellow))
	fmt.Println(box)
	fmt.Println(col("|                    QCLASS                     |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("| %-5d (%s)", q.QClass, classToString(q.QClass)), cYellow))
	fmt.Println(box)
}

// =============[ RESOURCE RECORDS ]=============

func printRRDiagram(index int, rr *DNSResourceRecord) {
	box := col("+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+--+", cDim)

	fmt.Println(box)
	fmt.Println(col(fmt.Sprintf(";; RR %d", index+1), cCyan+cBold))
	fmt.Println(col("|                     NAME                      |", cWhite))
	fmt.Printf("%s\n", col("| "+padRight(rr.Name, 45), cYellow))
	fmt.Println(box)
	fmt.Println(col("|      TYPE       |      CLASS      |    TTL    |", cWhite))
	fmt.Printf("%s\n", col(
		fmt.Sprintf("| %-5d (%-6s)| %-5d (%-6s)| %9d |",
			rr.Type, typeToString(rr.Type),
			rr.Class, classToString(rr.Class),
			rr.TTL),
		cYellow,
	))
	fmt.Println(box)
	fmt.Println(col("|                  RDLENGTH                     |", cWhite))
	fmt.Printf("%s\n", col(fmt.Sprintf("| %10d", rr.RDLength), cYellow))
	fmt.Println(box)
	fmt.Println(col("|                     RDATA                     |", cWhite))

	hexData := strings.ToUpper(hex.EncodeToString(rr.RData))
	if len(hexData) == 0 {
		fmt.Println(col("| (empty)", cDim))
	} else {
		for len(hexData) > 0 {
			chunk := hexData
			if len(chunk) > 32 {
				chunk = chunk[:32]
				hexData = hexData[32:]
			} else {
				hexData = ""
			}
			fmt.Printf("%s\n", col("| "+chunk, cGreen))
		}
	}
	fmt.Println(box)
}

// =============[ HELPERS ]=============

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func typeToString(t uint16) string {
	switch t {
	case 1:
		return "A"
	case 2:
		return "NS"
	case 5:
		return "CNAME"
	case 6:
		return "SOA"
	case 12:
		return "PTR"
	case 15:
		return "MX"
	case 16:
		return "TXT"
	case 28:
		return "AAAA"
	default:
		return "TYPE" + fmt.Sprint(t)
	}
}

func classToString(c uint16) string {
	switch c {
	case 1:
		return "IN"
	case 3:
		return "CH"
	case 4:
		return "HS"
	default:
		return "CLASS" + fmt.Sprint(c)
	}
}
