package backend

import (
	"github.com/gin-gonic/gin"
	"github.com/pboehm/ddns/shared"
	"log"
	"strings"
)

type Backend struct {
	config *shared.Config
	lookup *HostLookup
}

func NewBackend(config *shared.Config, lookup *HostLookup) *Backend {
	return &Backend{
		config: config,
		lookup: lookup,
	}
}

func (b *Backend) Run() error {
	r := gin.New()
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
