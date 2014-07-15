package main

import (
	"flag"
	"github.com/pboehm/ddns/connection"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	DdnsDomain string
)

func init() {
	flag.StringVar(&DdnsDomain, "domain", "",
		"The subdomain which should be handled by DDNS")
}

func ValidateCommandArgs() {
	if DdnsDomain == "" {
		log.Fatal("You have to supply the domain via --domain=DOMAIN")
	} else if !strings.HasPrefix(DdnsDomain, ".") {
		// get the domain in the right format
		DdnsDomain = "." + DdnsDomain
	}
}

func PrepareForExecution() string {
	flag.Parse()
	ValidateCommandArgs()

	if len(flag.Args()) != 1 {
		usage()
	}
	cmd := flag.Args()[0]

	return cmd
}

func main() {
	cmd = PrepareForExecution()

	conn := connection.OpenConnection()
	defer conn.Close()

	switch cmd {
	case "backend":
		log.Printf("Starting PDNS Backend\n")
		RunBackend(conn)
	case "web":
		log.Printf("Starting Web Service\n")
		RunWebService(conn)
	default:
		usage()
	}
}

func usage() {
	log.Fatal("Usage: ./ddns [backend|web]")
}
