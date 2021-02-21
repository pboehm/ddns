package shared

import (
	"github.com/jnovack/flag"
	"log"
	"os"
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
	CustomTemplatePath string

	fs *flag.FlagSet
}

func (c *Config) Initialize() {
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "DDNS", 0)
	fs.StringVar(&c.Domain, "domain", "",
		"The subdomain which should be handled by DDNS")

	fs.StringVar(&c.SOAFqdn, "soa_fqdn", "",
		"The FQDN of the DNS server which is returned as a SOA record")

	fs.StringVar(&c.ListenBackend, "listen-backend", ":8053",
		"Which socket should the backend web service use to bind itself")

	fs.StringVar(&c.ListenFrontend, "listen-frontend", ":8080",
		"Which socket should the frontend web service use to bind itself")

	fs.StringVar(&c.CustomTemplatePath, "custom-template-path", "",
		"A path to a custom `index.html` file that will be used instead of the default frontend template")

	fs.StringVar(&c.RedisHost, "redis-host", ":6379",
		"The Redis socket that should be used")

	fs.IntVar(&c.HostExpirationDays, "expiration-days", 10,
		"The number of days after a host is released when it is not updated")

	fs.BoolVar(&c.Verbose, "verbose", false,
		"Be more verbose")

	c.fs = fs
}

func (c *Config) Validate() {
	if err := c.fs.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing configuration: %v", err)
	}

	if c.Domain == "" {
		log.Fatal("You have to supply the domain via env variable DDNS_DOMAIN or command line flag --domain=DOMAIN")
	} else if !strings.HasPrefix(c.Domain, ".") {
		// get the domain in the right format
		c.Domain = "." + c.Domain
	}

	if c.SOAFqdn == "" {
		log.Fatal("You have to supply the server FQDN via env variable DDNS_SOA_FQDN or command line flag --soa_fqdn=FQDN")
	}
}
