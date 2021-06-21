package utils

import (
	"log"
	"net"
	"os"
)

func CreateUDPDNSSocket() (net.Conn, error) {
	dnsServer := net.JoinHostPort(DNSADDR, DNSPORT)
	return net.Dial("udp", dnsServer)
}

func GetInputDomains() []string {
	cmdArgs := os.Args

	if len(cmdArgs) < 2 {
		log.Fatalf("Usage %s <domain1> <domain2> ... <domain n>\n", cmdArgs[0])
	}

	return cmdArgs[1:]
}
