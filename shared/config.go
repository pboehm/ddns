package shared

import (
	"flag"
	"log"
	"strings"
)

type Config struct {
	Verbose            bool
	Domain             string
	SOAFqdn            string
	HostExpirationDays int
	ListenFrontend     string
	ListenBackend      string
	RedisHost          string
}

func (c *Config) Initialize() {
	flag.StringVar(&c.Domain, "domain", "",
		"The subdomain which should be handled by DDNS")

	flag.StringVar(&c.SOAFqdn, "soa_fqdn", "",
		"The FQDN of the DNS server which is returned as a SOA record")

	flag.StringVar(&c.ListenBackend, "listen-backend", ":8053",
		"Which socket should the backend web service use to bind itself")

	flag.StringVar(&c.ListenFrontend, "listen-frontend", ":8080",
		"Which socket should the frontend web service use to bind itself")

	flag.StringVar(&c.RedisHost, "redis", ":6379",
		"The Redis socket that should be used")

	flag.IntVar(&c.HostExpirationDays, "expiration-days", 10,
		"The number of days after a host is released when it is not updated")

	flag.BoolVar(&c.Verbose, "verbose", false,
		"Be more verbose")
}

func (c *Config) Validate() {
	flag.Parse()

	if c.Domain == "" {
		log.Fatal("You have to supply the domain via --domain=DOMAIN")
	} else if !strings.HasPrefix(c.Domain, ".") {
		// get the domain in the right format
		c.Domain = "." + c.Domain
	}

	if c.SOAFqdn == "" {
		log.Fatal("You have to supply the server FQDN via --soa_fqdn=FQDN")
	}
}
