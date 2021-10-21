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
		QuestionName:  encodeDomain(domain),
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

func NewDNSResourceHeader(name string, resourceType, class uint16,
	ttl uint32, length uint16) *DNSResourceHeader {
	return &DNSResourceHeader{
		Name:   name,
		Type:   resourceType,
		Class:  class,
		TTL:    ttl,
		Length: length,
	}
}

func NewDNSResource(header *DNSResourceHeader, data []byte) *DNSResource {
	return &DNSResource{
		Header: header,
		Data:   data,
	}
}

func NewDNSStringResponse(header *DNSHeader, answersLength int) *DNSStringResponse {
	return &DNSStringResponse{
		Header:  header,
		Answers: make([]*DNSStringAnswer, answersLength),
	}
}

func NewDNSStringAnswer(header *DNSResourceHeader, data string) *DNSStringAnswer {
	return &DNSStringAnswer{
		Header: header,
		Data:   data,
	}
}

func (query *DNSQuery) SendRequest(writter io.Writer) error {
	buffer := &bytes.Buffer{}

	binary.Write(buffer, binary.BigEndian, query.Header)
	binary.Write(buffer, binary.BigEndian, query.Question.QuestionName)
	binary.Write(buffer, binary.BigEndian, query.Question.QuestionType)
	binary.Write(buffer, binary.BigEndian, query.Question.QuestionClass)

	_, err := writter.Write(buffer.Bytes())
	return err
}

func GetResponse(reader io.Reader) []byte {
	response := make([]byte, MAX_RESPONSE_LENGTH)

	if _, err := reader.Read(response); err != nil {
		log.Println("Response Error:", err)
		return []byte{}
	}

	return response
}

func (response *DNSResponse) ToJSON(writter io.Writer) {
	if response.Header.QDCount != 1 {
		return
	}

	encoder := json.NewEncoder(writter)
	encodedResponse := NewDNSStringResponse(response.Header, len(response.Answers))

	for i, resource := range response.Answers {
		data := ""

		if resource.Header.Type == PTR {
			data = decodeDomain(resource.Data)
		} else if resource.Header.Type == AAAA || resource.Header.Type == A {
			ip := make(net.IP, len(resource.Data))
			copy(ip, resource.Data)
			data = ip.String()
		}

		stringRecord := NewDNSStringAnswer(resource.Header, data)
		encodedResponse.Answers[i] = stringRecord
	}

	if err := encoder.Encode(encodedResponse); err != nil {
		log.Fatal("Error at encoding the response:", err)
	}
}
