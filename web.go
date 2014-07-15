package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pboehm/ddns/connection"
	"html/template"
	"net"
	"net/http"
)

func RunWebService(conn *connection.RedisConnection) {

	// Add index template from bindata
	index_content, err := Asset("templates/index.html")
	HandleErr(err)

	html, err := template.New("index.html").Parse(string(index_content))
	HandleErr(err)

	r := gin.Default()
	r.HTMLTemplates = html

	r.GET("/", func(g *gin.Context) {
		g.HTML(200, "index.html", gin.H{"domain": DdnsDomain})
	})

	r.GET("/available/:hostname", func(c *gin.Context) {
		hostname := c.Params.ByName("hostname")

		c.JSON(200, gin.H{
			"available": !conn.HostExist(hostname),
		})
	})

	r.GET("/new/:hostname", func(c *gin.Context) {
		hostname := c.Params.ByName("hostname")

		if conn.HostExist(hostname) {
			c.String(403, "This hostname has already been registered.")
			return
		}

		host := &connection.Host{Hostname: hostname, Ip: "127.0.0.1"}
		host.GenerateAndSetToken()

		conn.SaveHost(host)

		c.JSON(200, gin.H{
			"hostname":    host.Hostname,
			"token":       host.Token,
			"update_link": fmt.Sprintf("/update/%s/%s", host.Hostname, host.Token),
		})
	})

	r.GET("/update/:hostname/:token", func(c *gin.Context) {
		hostname := c.Params.ByName("hostname")
		token := c.Params.ByName("token")

		if !conn.HostExist(hostname) {
			c.String(404,
				"This hostname has not been registered or is expired.")
			return
		}

		host := conn.GetHost(hostname)

		if host.Token != token {
			c.String(403,
				"You have supplied the wrong token to manipulate this host")
			return
		}

		ip, err := GetRemoteAddr(c.Req)
		if err != nil {
			c.String(500, "Your sender IP address is not in the right format")
		}

		host.Ip = ip
		conn.SaveHost(host)

		c.String(200, fmt.Sprintf("Your current IP is %s", ip))
	})

	r.Run(":8080")
}

// Get the Remote Address of the client. At First we try to get the
// X-Forwarded-For Header which holds the IP if we are behind a proxy,
// otherwise the RemoteAddr is used
func GetRemoteAddr(req *http.Request) (string, error) {
	header_data, ok := req.Header["X-Forwarded-For"]

	if ok {
		return header_data[0], nil
	} else {
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		return ip, err
	}
}
