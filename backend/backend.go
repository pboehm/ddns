package backend

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pboehm/ddns/config"
	"github.com/pboehm/ddns/hosts"
	"io"
	"strings"
	"time"
)

// PowerDnsBackend implements the PowerDNS-Pipe-Backend protocol ABI-Version 1
type PowerDnsBackend struct {
	config *config.Config
	hosts  hosts.HostBackend
	in     io.Reader
	out    io.Writer
}

// NewPowerDnsBackend creates an instance of a PowerDNS-Pipe-Backend using the supplied parameters
func NewPowerDnsBackend(config *config.Config, backend hosts.HostBackend, in io.Reader, out io.Writer) *PowerDnsBackend {
	return &PowerDnsBackend{
		config: config,
		hosts:  backend,
		in:     in,
		out:    out,
	}
}

// Run reads requests from an input (normally STDIN) and prints response messages to an output (normally STDOUT)
func (b *PowerDnsBackend) Run() {
	responses := make(chan backendResponse, 5)

	go func() {
		for response := range responses {
			fmt.Fprintln(b.out, strings.Join(response, "\t"))
		}
	}()

	// handshake with PowerDNS
	bio := bufio.NewReader(b.in)
	_, _, _ = bio.ReadLine()
	responses <- handshakeResponse

	for {
		request, err := b.parseRequest(bio)
		if err != nil {
			responses <- failResponse
			continue
		}

		if err = b.handleRequest(request, responses); err != nil {
			responses <- newResponse("LOG", err.Error())
		}
	}
}

// handleRequest handles the supplied request by sending response messages on the supplied responses channel
func (b *PowerDnsBackend) handleRequest(request *backendRequest, responses chan backendResponse) error {
	defer b.commitRequest(responses)

	responseRecord := request.queryType
	var response string

	switch request.queryType {
	case "SOA":
		response = fmt.Sprintf("%s. hostmaster%s. %d 1800 3600 7200 5",
			b.config.SOAFqdn, b.config.Domain, b.currentSOASerial())

	case "NS":
		response = fmt.Sprintf("%s.", b.config.SOAFqdn)

	case "A", "AAAA", "ANY":
		hostname, err := b.extractHostname(request.queryName)
		if err != nil {
			return err
		}

		var host *hosts.Host
		if host, err = b.hosts.GetHost(hostname); err != nil {
			return err
		}

		response = host.Ip

		responseRecord = "A"
		if !host.IsIPv4() {
			responseRecord = "AAAA"
		}

		if (request.queryType == "A" || request.queryType == "AAAA") && request.queryType != responseRecord {
			return errors.New("IP address is not valid for requested record")
		}

	default:
		return errors.New("Unsupported query type")
	}

	responses <- newResponse("DATA", request.queryName, request.queryClass, responseRecord, "10", request.queryId, response)

	return nil
}

func (b *PowerDnsBackend) commitRequest(responses chan backendResponse) {
	responses <- endResponse
}

// parseRequest reads a line from input and tries to build a request structure from it
func (b *PowerDnsBackend) parseRequest(input *bufio.Reader) (*backendRequest, error) {
	line, _, err := input.ReadLine()
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(line), "\t")
	if len(parts) != 6 {
		return nil, errors.New("Invalid line")
	}

	return &backendRequest{
		query:      parts[0],
		queryName:  parts[1],
		queryClass: parts[2],
		queryType:  parts[3],
		queryId:    parts[4],
	}, nil
}

// extractHostname extract the host part of the fqdn: pi.d.example.org -> pi
func (b *PowerDnsBackend) extractHostname(rawQueryName string) (string, error) {
	queryName := strings.ToLower(rawQueryName)

	hostname := ""
	if strings.HasSuffix(queryName, b.config.Domain) {
		hostname = queryName[:len(queryName)-len(b.config.Domain)]
	}

	if hostname == "" {
		return "", errors.New("Query name does not correspond to our domain")
	}

	return hostname, nil
}

// currentSOASerial get the current SOA serial by returning the current time in seconds
func (b *PowerDnsBackend) currentSOASerial() int64 {
	return time.Now().Unix()
}

type backendRequest struct {
	query      string
	queryName  string
	queryClass string
	queryType  string
	queryId    string
}

type backendResponse []string

func newResponse(values ...string) backendResponse {
	return values
}

var (
	handshakeResponse backendResponse = []string{"OK", "DDNS Backend"}
	endResponse       backendResponse = []string{"END"}
	failResponse      backendResponse = []string{"FAIL"}
)
