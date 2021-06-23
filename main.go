package main

import (
	"io"
	"log"
	"net"
	"os"
	"sync"

	"github.com/V-H-R-Oliveira/dns-client/protocol"
	"github.com/V-H-R-Oliveira/dns-client/utils"
)

func dnsQuery(wg *sync.WaitGroup, domain string, ipv6, reverseQuery bool, writter io.Writer) {
	defer wg.Done()

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
	query.SendRequest(socket)

	response := make([]byte, utils.MAX_RESPONSE_LENGTH)

	if _, err := socket.Read(response); err != nil {
		log.Println("Response Error:", err)
		return
	}

	_, res := protocol.ParseDNSResponse(response)
	res.ToJSON(writter)
}

func main() {
	inputs := utils.GetInputDomains()
	var wg sync.WaitGroup

	wg.Add(len(inputs))
	reverseQuery := false

	for _, domain := range inputs {
		ip := net.ParseIP(domain)

		if ip != nil {
			domain = utils.ReverseIPV4(ip)
			reverseQuery = true
		} else {
			reverseQuery = false
		}

		go dnsQuery(&wg, domain, false, reverseQuery, nil)
	}

	wg.Wait()
}
