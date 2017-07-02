package main

import (
	"github.com/pboehm/ddns/backend"
	"github.com/pboehm/ddns/config"
	"github.com/pboehm/ddns/hosts"
	"github.com/pboehm/ddns/web"
	"flag"
	"log"
	"os"
	"strings"
)

const (
	CmdBackend string = "backend"
	CmdWeb     string = "web"
)

var serviceConfig *config.Config = &config.Config{}

func init() {
	flag.StringVar(&serviceConfig.Domain, "domain", "",
		"The subdomain which should be handled by DDNS")

	flag.StringVar(&serviceConfig.Listen, "listen", ":8080",
		"Which socket should the web service use to bind itself")

	flag.StringVar(&serviceConfig.RedisHost, "redis", ":6379",
		"The Redis socket that should be used")

	flag.StringVar(&serviceConfig.SOAFqdn, "soa_fqdn", "",
		"The FQDN of the DNS server which is returned as a SOA record")

	flag.IntVar(&serviceConfig.HostExpirationDays, "expiration-days", 10,
		"The number of days after a host is released when it is not updated")

	flag.BoolVar(&serviceConfig.Verbose, "verbose", false,
		"Be more verbose")
}

func usage() {
	log.Fatal("Usage: ./ddns [backend|web]")
}

func validateCommandArgs(cmd string) {
	if serviceConfig.Domain == "" {
		log.Fatal("You have to supply the domain via --domain=DOMAIN")
	} else if !strings.HasPrefix(serviceConfig.Domain, ".") {
		// get the domain in the right format
		serviceConfig.Domain = "." + serviceConfig.Domain
	}

	if cmd == CmdBackend {
		if serviceConfig.SOAFqdn == "" {
			log.Fatal("You have to supply the server FQDN via --soa_fqdn=FQDN")
		}
	}
}

func prepareForExecution() string {
	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
	}
	cmd := flag.Args()[0]

	validateCommandArgs(cmd)
	return cmd
}

func main() {
	cmd := prepareForExecution()

	redis := hosts.NewRedisBackend(serviceConfig)
	defer redis.Close()

	switch cmd {
	case CmdBackend:
		backend.NewPowerDnsBackend(serviceConfig, redis, os.Stdin, os.Stdout).Run()

	case CmdWeb:
		web.NewWebService(serviceConfig, redis).Run()

	default:
		usage()
	}
}
