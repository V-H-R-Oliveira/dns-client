package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/V-H-R-Oliveira/dns-client/protocol"
	"github.com/V-H-R-Oliveira/dns-client/utils"
)

func dnsQuery(domain string, ipv6, reverseQuery bool, writter io.Writer) {
	if writter == nil {
		writter = os.Stdout
	}

	socket, err := utils.CreateUDPDNSSocket()

	if err != nil {
		log.Println("Error at creating the socket:", err)
		return
	}

	defer socket.Close()

	queryClass := protocol.A

	if ipv6 {
		queryClass = protocol.AAAA
	}

	if reverseQuery {
		queryClass = protocol.PTR
	}

	query := protocol.NewDNSQuery(domain, uint16(queryClass))

	if err := query.SendRequest(socket); err != nil {
		log.Fatalf("Failed to send a %s dns request due error: %s\n", domain, err.Error())
	}

	response := protocol.GetResponse(socket)
	_, res := protocol.ParseDNSResponse(response, false)

	res.ToJSON(writter)
}

func main() {
	inputs := utils.GetInputDomains()
	reverseQuery, ipv6 := false, false

	for _, domain := range inputs {
		ip := net.ParseIP(domain)

		if ip != nil {
			if ip.IsLoopback() || ip.IsUnspecified() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
				continue
			}

			if ip.DefaultMask() != nil {
				domain = utils.ReverseIPV4(ip)
			} else {
				domain = utils.ReverseIPV6(ip)
				ipv6 = true
			}

			reverseQuery = true
		}

		dnsQuery(domain, ipv6, reverseQuery, nil)
		ipv6, reverseQuery = false, false
	}
}
