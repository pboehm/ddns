package backend

import (
	"github.com/pboehm/ddns/shared"
	"gopkg.in/gin-gonic/gin.v1"
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
	r := gin.Default()

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
			log.Printf("Error during lookup: %v", err)
			c.JSON(200, gin.H{
				"result": false,
			})
		}
	})

	return r.Run(b.config.BackendListen)
}
