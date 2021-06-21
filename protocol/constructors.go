package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
)

func NewDNSHeader() *DNSHeader {
	return &DNSHeader{
		ID:      generateRequestID(),
		Flags:   RECDISERED | RECAVAILABLE,
		QDCount: 1,
	}
}

func NewDNSQuestion(domain string) *DNSQuestion {
	return &DNSQuestion{
		QuestionName:  EncodeDomain(domain),
		QuestionType:  A,
		QuestionClass: QCLASS,
	}
}

func NewDNSQuery(domain string) *DNSQuery {
	return &DNSQuery{
		Header:   NewDNSHeader(),
		Question: NewDNSQuestion(domain),
	}
}

func (query *DNSQuery) SendRequest(writter io.Writer) {
	buffer := &bytes.Buffer{}

	binary.Write(buffer, binary.BigEndian, query.Header)
	binary.Write(buffer, binary.BigEndian, query.Question.QuestionName)
	binary.Write(buffer, binary.BigEndian, query.Question.QuestionType)
	binary.Write(buffer, binary.BigEndian, query.Question.QuestionClass)

	writter.Write(buffer.Bytes())
}

func (response *DNSResponse) ToJSON(writter io.Writer) {
	encoder := json.NewEncoder(writter)

	if err := encoder.Encode(response); err != nil {
		log.Fatal("Error at encoding the response:", err)
	}
}
