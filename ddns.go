package main

import (
	"flag"
	"github.com/pboehm/ddns/backend"
	"github.com/pboehm/ddns/frontend"
	"github.com/pboehm/ddns/shared"
	"log"
	"strings"
)

var serviceConfig *shared.Config = &shared.Config{}

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

func validateCommandArgs() {
	if serviceConfig.Domain == "" {
		log.Fatal("You have to supply the domain via --domain=DOMAIN")
	} else if !strings.HasPrefix(serviceConfig.Domain, ".") {
		// get the domain in the right format
		serviceConfig.Domain = "." + serviceConfig.Domain
	}

	if serviceConfig.SOAFqdn == "" {
		log.Fatal("You have to supply the server FQDN via --soa_fqdn=FQDN")
	}
}

func main() {
	flag.Parse()
	validateCommandArgs()

	redis := shared.NewRedisBackend(serviceConfig)
	defer redis.Close()

	lookup := backend.NewHostLookup(serviceConfig, redis)

	frontend.NewWebService(serviceConfig, redis, lookup).Run()
}
