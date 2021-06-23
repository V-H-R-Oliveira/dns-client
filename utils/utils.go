package utils

import (
	"log"
	"net"
	"os"

	"github.com/V-H-R-Oliveira/dns-client/protocol"
)

func CreateUDPDNSSocket() (net.Conn, error) {
	dnsServer := net.JoinHostPort(DNS_ADDR, DNS_PORT)
	return net.Dial("udp", dnsServer)
}

func GetInputDomains() []string {
	cmdArgs := os.Args

	if len(cmdArgs) < 2 {
		log.Fatalf("Usage %s <domain1> <domain2> ... <domain n>\n", cmdArgs[0])
	}

	return cmdArgs[1:]
}

func ReverseIPV4(ip net.IP) string {
	reversedIp := net.IP{}

	for i := 0; i < 4; i++ {
		octet := ip[len(ip)-1-i]
		reversedIp = append(reversedIp, octet)
	}

	return reversedIp.To4().String() + protocol.REVERSE_DNS_DOMAIN
}
