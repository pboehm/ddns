package main

import (
	"bufio"
	"fmt"
	"github.com/pboehm/ddns/connection"
	"os"
	"strings"
)

// This function implements the PowerDNS-Pipe-Backend protocol and generates
// the response data it possible
func RunBackend(conn *connection.RedisConnection) {

	// handshake with PowerDNS
	fmt.Printf("OK\tDDNS Go Backend\n")

	bio := bufio.NewReader(os.Stdin)

	for {
		line, _, err := bio.ReadLine()

		HandleErr(err)
		HandleRequest(string(line), conn)
	}
}

func HandleRequest(line string, conn *connection.RedisConnection) {
	defer fmt.Printf("END\n")

	parts := strings.Split(line, "\t")
	if len(parts) != 6 {
		return
	}

	query_name := parts[1]
	query_class := parts[2]
	// query_type  := parts[3] // TODO Handle SOA Requests
	query_id := parts[4]

	// get the host part of the fqdn
	// pi.d.example.org -> pi
	hostname := ""
	if strings.HasSuffix(query_name, DdnsDomain) {
		hostname = query_name[:len(query_name)-len(DdnsDomain)]
	}

	if hostname == "" || !conn.HostExist(hostname) {
		return
	}

	host := conn.GetHost(hostname)

	record := "A"
	if !host.IsIPv4() {
		record = "AAAA"
	}

	fmt.Printf("DATA\t%s\t%s\t%s\t10\t%s\t%s\n",
		query_name, query_class, record, query_id, host.Ip)
}
