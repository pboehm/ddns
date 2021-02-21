package backend

import (
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/pboehm/ddns/shared"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"strings"
)

type Backend struct {
	config   *shared.Config
	lookup   *HostLookup
	registry *prometheus.Registry
}

func NewBackend(config *shared.Config, lookup *HostLookup, registry *prometheus.Registry) *Backend {
	return &Backend{
		config:   config,
		lookup:   lookup,
		registry: registry,
	}
}

func (b *Backend) Run() error {
	prom := ginprom.New(ginprom.Namespace("ddns"),
		ginprom.Subsystem("backend"), ginprom.Registry(b.registry))

	r := gin.New()
	r.Use(prom.Instrument())
	prom.Use(r)

	r.Use(gin.Recovery())

	if b.config.Verbose {
		r.Use(gin.Logger())
	}

	r.GET("/dnsapi/lookup/:qname/:qtype", func(c *gin.Context) {
		request := &Request{
			QName: strings.TrimRight(c.Param("qname"), "."),
			QType: c.Param("qtype"),
		}

		response, err := b.lookup.Lookup(request)
		if err == nil {
			c.JSON(200, gin.H{
				"result": []*Response{response},
			})
		} else {
			if b.config.Verbose {
				log.Printf("Error during lookup: %v", err)
			}

			c.JSON(200, gin.H{
				"result": false,
			})
		}
	})

	r.GET("/dnsapi/getDomainMetadata/:name/:kind", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"result": false,
		})
	})

	return r.Run(b.config.ListenBackend)
}
