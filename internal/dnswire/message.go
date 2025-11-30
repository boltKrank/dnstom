package dnswire

import (
	"fmt"
	"io"
)

type Header struct {
	ID      uint16
	QR      bool
	Opcode  uint8
	AA      bool
	TC      bool
	RD      bool
	RA      bool
	Z       uint8 // usually 0
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
}

type ResourceRecord struct {
	Name     string
	Type     uint16
	Class    uint16
	TTL      uint32
	RDLength uint16
	RData    []byte // you'll later decode into typed RRs using internal/rr
}

type Message struct {
	Header     Header
	Questions  []Question
	Answers    []ResourceRecord
	Authority  []ResourceRecord
	Additional []ResourceRecord
}

// PrettyPrint is a placeholder for now; you'll make this nicer later.
func PrettyPrint(m *Message, w io.Writer) error {
	_, err := fmt.Fprintf(w, ";; dnstom raw message\n;; ID: %d, QD: %d, AN: %d, NS: %d, AR: %d\n",
		m.Header.ID, m.Header.QDCount, m.Header.ANCount, m.Header.NSCount, m.Header.ARCount)
	if err != nil {
		return err
	}

	for _, q := range m.Questions {
		if _, err := fmt.Fprintf(w, ";; QUESTION: %s type=%d class=%d\n", q.Name, q.Type, q.Class); err != nil {
			return err
		}
	}

	for _, rr := range m.Answers {
		if _, err := fmt.Fprintf(w, "ANSWER: %s TTL=%d type=%d class=%d\n",
			rr.Name, rr.TTL, rr.Type, rr.Class); err != nil {
			return err
		}
	}

	return nil
}
