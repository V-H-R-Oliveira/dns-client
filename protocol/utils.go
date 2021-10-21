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
		log.Fatal("Failed to generate a random request id due error:", err)
	}

	return binary.BigEndian.Uint16(buffer)
}

func encodeDomain(domain string) []byte {
	cleanedDomain := strings.Trim(domain, "\n\r\t ")
	splittedDomain := strings.Split(cleanedDomain, ".")
	s := ""

	for _, slice := range splittedDomain {
		label := fmt.Sprintf("%c%s", len(slice), slice)

		if len(label) > MAX_LABEL_LENGTH {
			continue
		}

		s += label
	}

	s += fmt.Sprintf("%c", 0)
	return []byte(s)
}

func decodeDomain(domain []byte) string {
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

func logResponseStatus(id, flag uint16) {
	switch flag & 0xf {
	case RES_OK:
		log.Printf("[0x%04x] Response Status: OK\n", id)
		return
	case RES_FORMAT_ERROR:
		log.Printf("[0x%04x] Response Status: FORMAT ERROR\n", id)
		return
	case RES_SERVER_FAILURE:
		log.Printf("[0x%04x] Response Status: SERVER FAILURE\n", id)
		return
	case RES_NAME_ERROR:
		log.Printf("[0x%04x] Response Status: NAME ERROR\n", id)
		return
	case RES_NOT_IMPLEMENTED:
		log.Printf("[0x%04x] Response Status: NOT IMPLEMENTED\n", id)
		return
	case RES_REFUSED:
		log.Printf("[0x%04x] Response Status: REFUSED\n", id)
		return
	default:
		log.Printf("[0x%04x] Response Status: UNKNOWN\n", id)
	}
}

func parseHeader(header []byte, debug bool) *DNSHeader {
	id := binary.BigEndian.Uint16(header[:2])
	flags := binary.BigEndian.Uint16(header[2:4])

	if debug {
		logResponseStatus(id, flags)
	}

	return &DNSHeader{
		ID:      id,
		Flags:   flags,
		QDCount: binary.BigEndian.Uint16(header[4:6]),
		ANCount: binary.BigEndian.Uint16(header[6:8]),
		NSCount: binary.BigEndian.Uint16(header[8:10]),
		ARCount: binary.BigEndian.Uint16(header[10:]),
	}
}

func parseQuestion(question []byte) (*DNSQuestion, int) {
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

	return decodeDomain(response[:endDomainOffset])
}

func parseAnswer(fullResponse, answers []byte, resourcesAmount uint16) []*DNSResource {
	if resourcesAmount == 0 {
		return []*DNSResource{}
	}

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

		resourceType := binary.BigEndian.Uint16(frame[1:3])
		resourceClass := binary.BigEndian.Uint16(frame[3:5])
		ttl := binary.BigEndian.Uint32(frame[5:9])
		length := binary.BigEndian.Uint16(frame[9:11])
		data := frame[11:]

		resourceHeader := NewDNSResourceHeader(domain, resourceType, resourceClass, ttl, length)
		resources[i] = NewDNSResource(resourceHeader, data)
	}

	return resources
}

func ParseDNSResponse(response []byte, debug bool) (*DNSQuery, *DNSResponse) {
	responsePayload := bytes.Trim(response, "\n\r\t\x00")
	fullResponseCopy := responsePayload
	header := responsePayload[:12]

	dnsHeader := parseHeader(header, debug)
	responsePayload = responsePayload[12:]

	dnsQuestion, startAnswerOffset := parseQuestion(responsePayload)
	responsePayload = responsePayload[startAnswerOffset:]

	dnsAnswers := parseAnswer(fullResponseCopy, responsePayload, dnsHeader.ANCount)

	return &DNSQuery{
			Header:   dnsHeader,
			Question: dnsQuestion,
		}, &DNSResponse{
			Header:  dnsHeader,
			Answers: dnsAnswers,
		}
}
