package main

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/V-H-R-Oliveira/dns-client/protocol"
	"github.com/V-H-R-Oliveira/dns-client/utils"
)

func worker(wg *sync.WaitGroup, domain string, writter io.Writer) {
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

	query := protocol.NewDNSQuery(domain)
	query.SendRequest(socket)

	response := make([]byte, utils.MAXLENGTH)

	if _, err := socket.Read(response); err != nil {
		log.Println("Response Error:", err)
		return
	}

	_, res := protocol.ParseDNSResponse(response)
	res.ToJSON(writter)
}

// TODO: Add support to Reverse DNS queries
func main() {
	inputs := utils.GetInputDomains()
	var wg sync.WaitGroup

	wg.Add(len(inputs))

	for _, domain := range inputs {
		go worker(&wg, domain, nil)
	}

	wg.Wait()
}
