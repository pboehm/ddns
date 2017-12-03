package main

import (
	"flag"
	"github.com/pboehm/ddns/backend"
	"github.com/pboehm/ddns/frontend"
	"github.com/pboehm/ddns/shared"
	"golang.org/x/sync/errgroup"
	"log"
	"strings"
)

var serviceConfig *shared.Config = &shared.Config{}

func init() {
	flag.StringVar(&serviceConfig.Domain, "domain", "",
		"The subdomain which should be handled by DDNS")

	flag.StringVar(&serviceConfig.SOAFqdn, "soa_fqdn", "",
		"The FQDN of the DNS server which is returned as a SOA record")

	flag.StringVar(&serviceConfig.BackendListen, "listen-backend", ":8057",
		"Which socket should the backend web service use to bind itself")

	flag.StringVar(&serviceConfig.FrontendListen, "listen-frontend", ":8080",
		"Which socket should the frontend web service use to bind itself")

	flag.StringVar(&serviceConfig.RedisHost, "redis", ":6379",
		"The Redis socket that should be used")

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

	var group errgroup.Group

	group.Go(func() error {
		lookup := backend.NewHostLookup(serviceConfig, redis)
		return backend.NewBackend(serviceConfig, lookup).Run()
	})

	group.Go(func() error {
		return frontend.NewFrontend(serviceConfig, redis).Run()
	})

	if err := group.Wait(); err != nil {
		log.Fatal(err)
	}
}
