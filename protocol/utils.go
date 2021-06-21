package protocol

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)

func generateRequestID() uint16 {
	buffer := make([]byte, 2)

	if _, err := rand.Read(buffer); err != nil {
		log.Fatal(err)
	}

	return binary.BigEndian.Uint16(buffer)
}

func EncodeDomain(domain string) []byte {
	cleanedDomain := strings.Trim(domain, "\n\r\t ")
	splittedDomain := strings.Split(cleanedDomain, ".")
	s := ""

	for _, slice := range splittedDomain {
		s += fmt.Sprintf("%c%s", len(slice), slice)
	}

	s += fmt.Sprintf("%c", 0)
	return []byte(s)
}

func DecodeDomain(domain []byte) string {
	decodedDomain := []byte{}
	limit, counter := 0, 0

	for i, element := range domain {
		if i == 0 {
			limit = int(element)
			continue
		}

		if counter == limit {
			limit = int(element)
			decodedDomain = append(decodedDomain, '.')
			counter = 0
			continue
		}

		decodedDomain = append(decodedDomain, element)
		counter++
	}

	return string(decodedDomain)
}

func ParseHeader(header []byte) *DNSHeader {
	return &DNSHeader{
		ID:      binary.BigEndian.Uint16(header[:2]),
		Flags:   binary.BigEndian.Uint16(header[2:4]),
		QDCount: binary.BigEndian.Uint16(header[4:6]),
		ANCount: binary.BigEndian.Uint16(header[6:8]),
		NSCount: binary.BigEndian.Uint16(header[8:10]),
		ARCount: binary.BigEndian.Uint16(header[10:]),
	}
}

func ParseQuestion(question []byte) (*DNSQuestion, int) {
	domainEnd := bytes.IndexByte(question, 0)

	if domainEnd == -1 {
		log.Fatal("Failed to find the question domain start offset")
	}

	questionName := question[:domainEnd]
	startAnswerOffset := domainEnd + 5
	questionTypeClass := question[domainEnd+1 : startAnswerOffset] // skip 0x00 marker

	return &DNSQuestion{
		QuestionName:  questionName,
		QuestionType:  binary.BigEndian.Uint16(questionTypeClass[:2]),
		QuestionClass: binary.BigEndian.Uint16(questionTypeClass[2:4]),
	}, startAnswerOffset
}

func fetchDomainFromResponse(response []byte, offset int) string {
	response = response[offset:]
	endDomainOffset := bytes.IndexByte(response, 0)

	if endDomainOffset == -1 {
		log.Fatal("Failed to find the question domain end offset")
	}

	return DecodeDomain(response[:endDomainOffset])
}

func ParseAnswer(fullResponse, answers []byte, resourcesAmount uint16) []*DNSResource {
	offset := answers[0]
	answers = answers[1:]
	answersFrames := bytes.Split(answers, []byte{offset})
	resources := make([]*DNSResource, resourcesAmount)
	cache := make(map[int]string)

	for i, frame := range answersFrames {
		domainOffset := int(frame[0])
		domain, ok := cache[domainOffset]

		if !ok {
			domain = fetchDomainFromResponse(fullResponse, domainOffset)
			cache[domainOffset] = domain
		}

		resources[i] = &DNSResource{
			Name:   domain,
			Type:   binary.BigEndian.Uint16(frame[1:3]),
			Class:  binary.BigEndian.Uint16(frame[3:5]),
			TTL:    binary.BigEndian.Uint32(frame[5:9]),
			Length: binary.BigEndian.Uint16(frame[9:11]),
			Data:   frame[11:],
		}
	}

	return resources
}

func ParseDNSResponse(response []byte) (*DNSQuery, *DNSResponse) {
	responsePayload := bytes.Trim(response, "\n\r\t\x00")
	fullResponseCopy := responsePayload
	header := responsePayload[:12]

	dnsHeader := ParseHeader(header)
	responsePayload = responsePayload[12:]

	dnsQuestion, startAnswerOffset := ParseQuestion(responsePayload)
	responsePayload = responsePayload[startAnswerOffset:]

	dnsAnswers := ParseAnswer(fullResponseCopy, responsePayload, dnsHeader.ANCount)

	return &DNSQuery{
			Header:   dnsHeader,
			Question: dnsQuestion,
		}, &DNSResponse{
			Header:  dnsHeader,
			Answers: dnsAnswers,
		}
}
