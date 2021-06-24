package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/V-H-R-Oliveira/dns-client/protocol"
	"github.com/joho/godotenv"
)

func CreateUDPDNSSocket() (net.Conn, error) {
	dnsServer := net.JoinHostPort(DNS_ADDR, DNS_PORT)
	return net.Dial("udp", dnsServer)
}

func GetInputDomains() []string {
	cmdArgs := os.Args

	if len(cmdArgs) < 2 {
		log.Fatalf("Usage %s <domain 1 | ip 1> <domain 2 | ip 2> ... <domain n | ip n>\nExamples: %s google.com microsoft.com 8.8.4.4\n", cmdArgs[0], cmdArgs[0])
	}

	return cmdArgs[1:]
}

func ReverseIPV4(ip net.IP) string {
	reversedIp := net.IP{}

	for i := 0; i < 4; i++ {
		octet := ip[len(ip)-1-i]
		reversedIp = append(reversedIp, octet)
	}

	return reversedIp.To4().String() + protocol.REVERSE_DNS_IPV4_DOMAIN
}

func ReverseIPV6(ip net.IP) string {
	reversedIp := ""

	for i := 0; i < 16; i++ {
		nibble := ip[len(ip)-1-i]
		reversedIp += fmt.Sprintf("%x.%x.", nibble&0xf, nibble>>4)
	}

	return reversedIp + protocol.REVERSE_DNS_IPV6_DOMAIN
}

func LoadEnv() {
	godotenv.Load()
}

func IsDebugMode() bool {
	debugEnvVar, ok := os.LookupEnv("DEBUG")

	if !ok || debugEnvVar == "" {
		return false
	}

	debug, err := strconv.ParseBool(debugEnvVar)

	if err != nil {
		log.Fatal("DEBUG variable should be a valid boolean value.")
	}

	return debug
}
