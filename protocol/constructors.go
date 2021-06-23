package protocol

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"
)

func NewDNSHeader() *DNSHeader {
	return &DNSHeader{
		ID:      generateRequestID(),
		Flags:   REC_DISERED | REC_AVAILABLE,
		QDCount: 1,
	}
}

func NewDNSQuestion(domain string, class uint16) *DNSQuestion {
	return &DNSQuestion{
		QuestionName:  EncodeDomain(domain),
		QuestionType:  class,
		QuestionClass: QCLASS,
	}
}

func NewDNSQuery(domain string, class uint16) *DNSQuery {
	return &DNSQuery{
		Header:   NewDNSHeader(),
		Question: NewDNSQuestion(domain, class),
	}
}

func NewDNSStringResponse(header *DNSHeader, answersLength int) *DNSStringResponse {
	return &DNSStringResponse{
		Header:  header,
		Answers: make([]*DNSStringAnswer, answersLength),
	}
}

func NewDNSStringAnswer(name string, resourceType, class uint16,
	ttl uint32, length uint16, data string) *DNSStringAnswer {
	return &DNSStringAnswer{
		Name:   name,
		Type:   resourceType,
		Class:  class,
		TTL:    ttl,
		Length: length,
		Data:   data,
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
	if response.Header.QDCount != 1 {
		return
	}

	encoder := json.NewEncoder(writter)
	encodedResponse := NewDNSStringResponse(response.Header, len(response.Answers))

	for i, resource := range response.Answers {
		data := ""

		if resource.Type == PTR {
			data = DecodeDomain(resource.Data)
		} else if resource.Type == AAAA || resource.Type == A {
			ip := make(net.IP, len(resource.Data))
			copy(ip, resource.Data)
			data = ip.String()
		}

		stringRecord := NewDNSStringAnswer(resource.Name, resource.Type,
			resource.Class, resource.TTL, resource.Length, data)

		encodedResponse.Answers[i] = stringRecord
	}

	if err := encoder.Encode(encodedResponse); err != nil {
		log.Fatal("Error at encoding the response:", err)
	}
}
