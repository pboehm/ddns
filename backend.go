package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type responder func()

func respondWithFAIL() {
	fmt.Printf("FAIL\n")
}

func respondWithEND() {
	fmt.Printf("END\n")
}

// This function implements the PowerDNS-Pipe-Backend protocol and generates
// the response data it possible
func RunBackend(conn *RedisConnection) {
	bio := bufio.NewReader(os.Stdin)

	// handshake with PowerDNS
	_, _, _ = bio.ReadLine()
	fmt.Printf("OK\tDDNS Go Backend\n")

	for {
		line, _, err := bio.ReadLine()
		if err != nil {
			respondWithFAIL()
			continue
		}

		HandleRequest(string(line), conn)()
	}
}

func HandleRequest(line string, conn *RedisConnection) responder {
	if Verbose {
		fmt.Printf("LOG\t'%s'\n", line)
	}

	parts := strings.Split(line, "\t")
	if len(parts) != 6 {
		return respondWithFAIL
	}

	query_name := parts[1]
	query_class := parts[2]
	query_type := parts[3]
	query_id := parts[4]

	var response, record string
	record = query_type

	switch query_type {
	case "SOA":
		response = fmt.Sprintf("%s. hostmaster.example.com. %d 1800 3600 7200 5",
			DdnsSoaFqdn, getSoaSerial())

	case "NS":
		response = DdnsSoaFqdn

	case "A":
	case "ANY":
		// get the host part of the fqdn: pi.d.example.org -> pi
		hostname := ""
		if strings.HasSuffix(query_name, DdnsDomain) {
			hostname = query_name[:len(query_name)-len(DdnsDomain)]
		}

		if hostname == "" || !conn.HostExist(hostname) {
			return respondWithFAIL
		}

		host := conn.GetHost(hostname)
		response = host.Ip

		record = "A"
		if !host.IsIPv4() {
			record = "AAAA"
		}

	default:
		return respondWithFAIL
	}

	fmt.Printf("DATA\t%s\t%s\t%s\t10\t%s\t%s\n",
		query_name, query_class, record, query_id, response)
	return respondWithEND
}

func getSoaSerial() int64 {
	// return current time in milliseconds
	return time.Now().UnixNano()
}
