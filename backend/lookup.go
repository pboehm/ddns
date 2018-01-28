package backend

import (
	"errors"
	"fmt"
	"github.com/pboehm/ddns/shared"
	"strings"
	"time"
)

type Request struct {
	QType      string
	QName      string
	Remote     string
	Local      string
	RealRemote string
	ZoneId     string
}

type Response struct {
	QType   string `json:"qtype"`
	QName   string `json:"qname"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type HostLookup struct {
	config *shared.Config
	hosts  shared.HostBackend
}

func NewHostLookup(config *shared.Config, hostsBackend shared.HostBackend) *HostLookup {
	return &HostLookup{config, hostsBackend}
}

func (l *HostLookup) Lookup(request *Request) (*Response, error) {
	responseRecord := request.QType
	responseContent := ""

	switch request.QType {
	case "SOA":
		responseContent = fmt.Sprintf("%s. hostmaster%s. %d 1800 3600 7200 5",
			l.config.SOAFqdn, l.config.Domain, l.currentSOASerial())

	case "NS":
		responseContent = l.config.SOAFqdn

	case "A", "AAAA", "ANY":
		hostname, err := l.extractHostname(request.QName)
		if err != nil {
			return nil, err
		}

		var host *shared.Host
		if host, err = l.hosts.GetHost(hostname); err != nil {
			return nil, err
		}

		responseContent = host.Ip

		responseRecord = "A"
		if !host.IsIPv4() {
			responseRecord = "AAAA"
		}

		if (request.QType == "A" || request.QType == "AAAA") && request.QType != responseRecord {
			return nil, errors.New("IP address is not valid for requested record")
		}

	default:
		return nil, errors.New("Invalid request")
	}

	return &Response{QType: responseRecord, QName: request.QName, Content: responseContent, TTL: 5}, nil
}

// extractHostname extract the host part of the fqdn: pi.d.example.org -> pi
func (l *HostLookup) extractHostname(rawQueryName string) (string, error) {
	queryName := strings.ToLower(rawQueryName)

	hostname := ""
	if strings.HasSuffix(queryName, l.config.Domain) {
		hostname = queryName[:len(queryName)-len(l.config.Domain)]
	}

	if hostname == "" {
		return "", errors.New("Query name does not correspond to our domain")
	}

	return hostname, nil
}

// currentSOASerial get the current SOA serial by returning the current time in seconds
func (l *HostLookup) currentSOASerial() int64 {
	return time.Now().Unix()
}
