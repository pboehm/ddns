package main

import (
	"flag"
	"log"
	"strings"
)

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

const (
	CmdBackend string = "backend"
	CmdWeb     string = "web"
)

var (
	DdnsDomain          string
	DdnsWebListenSocket string
	DdnsRedisHost       string
	DdnsSoaFqdn         string
	Verbose             bool
)

func init() {
	flag.StringVar(&DdnsDomain, "domain", "",
		"The subdomain which should be handled by DDNS")

	flag.StringVar(&DdnsWebListenSocket, "listen", ":8080",
		"Which socket should the web service use to bind itself")

	flag.StringVar(&DdnsRedisHost, "redis", ":6379",
		"The Redis socket that should be used")

	flag.StringVar(&DdnsSoaFqdn, "soa_fqdn", "",
		"The FQDN of the DNS server which is returned as a SOA record")

	flag.BoolVar(&Verbose, "verbose", false,
		"Be more verbose")
}

func ValidateCommandArgs(cmd string) {
	if DdnsDomain == "" {
		log.Fatal("You have to supply the domain via --domain=DOMAIN")
	} else if !strings.HasPrefix(DdnsDomain, ".") {
		// get the domain in the right format
		DdnsDomain = "." + DdnsDomain
	}

	if cmd == CmdBackend {
		if DdnsSoaFqdn == "" {
			log.Fatal("You have to supply the server FQDN via --soa_fqdn=FQDN")
		}
	}
}

func PrepareForExecution() string {
	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
	}
	cmd := flag.Args()[0]

	ValidateCommandArgs(cmd)
	return cmd
}

func main() {
	cmd := PrepareForExecution()

	conn := OpenConnection(DdnsRedisHost)
	defer conn.Close()

	switch cmd {
	case CmdBackend:
		log.Printf("Starting PDNS Backend\n")
		RunBackend(conn)
	case CmdWeb:
		log.Printf("Starting Web Service\n")
		RunWebService(conn)
	default:
		usage()
	}
}

func usage() {
	log.Fatal("Usage: ./ddns [backend|web]")
}
