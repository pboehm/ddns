package main

import (
	"bufio"
	"fmt"
	"github.com/pboehm/ddns/connection"
	"os"
	"strings"
)

func RunBackend(conn *connection.RedisConnection) {

	fmt.Printf("OK\tDDNS Go Backend\n")

	bio := bufio.NewReader(os.Stdin)

	for {
		line, _, err := bio.ReadLine()
		HandleErr(err)

		parts := strings.Split(string(line), "\t")
		if len(parts) == 6 {
			query_name := parts[1]

			// get the host part of the fqdn
			// pi.d.example.org -> pi
			hostname := ""
			if strings.HasSuffix(query_name, DdnsDomain) {
				hostname = query_name[:len(query_name)-len(DdnsDomain)]
			}

			query_class := parts[2]
			// query_type := parts[3]
			query_id := parts[4]

			if hostname != "" {
				// check for existance and create response
				if conn.HostExist(hostname) {
					host := conn.GetHost(hostname)

					record := "A"
					if !host.IsIPv4() {
						record = "AAAA"
					}

					fmt.Printf("DATA\t%s\t%s\t%s\t10\t%s\t%s\n",
						query_name, query_class, record, query_id, host.Ip)
				}
			}
		}

		fmt.Printf("END\n")
	}
}
