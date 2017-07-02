package backend

import (
	"bufio"
	"bytes"
	"errors"
	c "github.com/pboehm/ddns/config"
	h "github.com/pboehm/ddns/hosts"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type testHostBackend struct {
	hosts map[string]*h.Host
}

func (b *testHostBackend) GetHost(hostname string) (*h.Host, error) {
	host, ok := b.hosts[hostname]
	if ok {
		return host, nil
	} else {
		return nil, errors.New("Host not found")
	}
}

func (b *testHostBackend) SetHost(host *h.Host) error {
	b.hosts[host.Hostname] = host
	return nil
}

func buildBackend(domain string) (*c.Config, *testHostBackend, *PowerDnsBackend) {
	config := &c.Config{
		Verbose: false,
		Domain:  domain,
		SOAFqdn: "dns" + domain,
	}

	hosts := &testHostBackend{
		hosts: map[string]*h.Host{
			"www": {
				Hostname: "www",
				Ip:       "10.11.12.13",
				Token:    "abcdef",
			},
			"v4": {
				Hostname: "v4",
				Ip:       "10.10.10.10",
				Token:    "ghijkl",
			},
			"v6": {
				Hostname: "v6",
				Ip:       "2001:db8:85a3::8a2e:370:7334",
				Token:    "ghijkl",
			},
		},
	}

	return config, hosts, NewPowerDnsBackend(config, hosts, os.Stdin, os.Stdout)
}

func buildRequest(queryName, queryType string) *backendRequest {
	return &backendRequest{
		query:      "Q",
		queryName:  queryName,
		queryClass: "IN",
		queryType:  queryType,
		queryId:    "-1",
	}
}

func readResponse(t *testing.T, responses chan backendResponse) backendResponse {
	select {
	case response, ok := <-responses:
		assert.True(t, ok)
		return response
	default:
		assert.FailNow(t, "Couldn't read response because it is not available ...")
		return nil
	}
}

func TestParseRequest(t *testing.T) {
	_, _, backend := buildBackend(".example.org")

	reader := bufio.NewReader(bytes.NewBufferString("Q\twww.example.org\tIN\tCNAME\t-1\t203.0.113.210\n"))
	request, err := backend.parseRequest(reader)
	assert.Nil(t, err)
	assert.Equal(t, buildRequest("www.example.org", "CNAME"), request)

	reader = bufio.NewReader(bytes.NewBufferString("Q\texample.org\tIN\tSOA\t-1\t203.0.113.210\n"))
	request, err = backend.parseRequest(reader)
	assert.Nil(t, err)
	assert.Equal(t, buildRequest("example.org", "SOA"), request)

	reader = bufio.NewReader(bytes.NewBufferString("Q\texample.org\n"))
	request, err = backend.parseRequest(reader)
	assert.NotNil(t, err)
	assert.Nil(t, request)
}

func TestRequestHandling(t *testing.T) {
	_, _, backend := buildBackend(".example.org")

	responses := make(chan backendResponse, 2)
	err := backend.handleRequest(buildRequest("example.org", "SOA"), responses)
	assert.Nil(t, err)
	soaResponse := readResponse(t, responses)
	assert.Equal(t, newResponse("DATA", "example.org", "IN", "SOA", "10", "-1"), soaResponse[0:6])
	assert.Regexp(t, "dns\\.example\\.org\\. hostmaster\\.example.org\\. \\d+ 1800 3600 7200 5", soaResponse[6])
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("example.org", "NS"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "example.org", "IN", "NS", "10", "-1", "dns.example.org."), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("www.example.org", "ANY"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "www.example.org", "IN", "A", "10", "-1", "10.11.12.13"), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("www.example.org", "A"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "www.example.org", "IN", "A", "10", "-1", "10.11.12.13"), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	// Allow hostname to be mixed case which is used by Let's Encrypt for a little bit more security
	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("wWW.eXaMPlE.oRg", "A"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "wWW.eXaMPlE.oRg", "IN", "A", "10", "-1", "10.11.12.13"), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 1)
	err = backend.handleRequest(buildRequest("notexisting.example.org", "A"), responses)
	assert.NotNil(t, err)
	assert.Equal(t, endResponse, readResponse(t, responses))

	// Correct Handling of IPv4/IPv6 and ANY/A/AAAA
	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("v4.example.org", "ANY"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "v4.example.org", "IN", "A", "10", "-1", "10.10.10.10"), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("v4.example.org", "A"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "v4.example.org", "IN", "A", "10", "-1", "10.10.10.10"), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 1)
	err = backend.handleRequest(buildRequest("v4.example.org", "AAAA"), responses)
	assert.NotNil(t, err)
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("v6.example.org", "ANY"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "v6.example.org", "IN", "AAAA", "10", "-1", "2001:db8:85a3::8a2e:370:7334"), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 2)
	err = backend.handleRequest(buildRequest("v6.example.org", "AAAA"), responses)
	assert.Nil(t, err)
	assert.Equal(t, newResponse("DATA", "v6.example.org", "IN", "AAAA", "10", "-1", "2001:db8:85a3::8a2e:370:7334"), readResponse(t, responses))
	assert.Equal(t, endResponse, readResponse(t, responses))

	responses = make(chan backendResponse, 1)
	err = backend.handleRequest(buildRequest("v6.example.org", "A"), responses)
	assert.NotNil(t, err)
	assert.Equal(t, endResponse, readResponse(t, responses))
}
