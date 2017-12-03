package backend

import (
	"errors"
	"github.com/pboehm/ddns/shared"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testHostBackend struct {
	hosts map[string]*shared.Host
}

func (b *testHostBackend) GetHost(hostname string) (*shared.Host, error) {
	host, ok := b.hosts[hostname]
	if ok {
		return host, nil
	} else {
		return nil, errors.New("Host not found")
	}
}

func (b *testHostBackend) SetHost(host *shared.Host) error {
	b.hosts[host.Hostname] = host
	return nil
}

func buildLookup(domain string) (*shared.Config, *testHostBackend, *HostLookup) {
	config := &shared.Config{
		Verbose: false,
		Domain:  domain,
		SOAFqdn: "dns" + domain,
	}

	hosts := &testHostBackend{
		hosts: map[string]*shared.Host{
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

	return config, hosts, &HostLookup{config, hosts}
}

func buildRequest(queryName, queryType string) *Request {
	return &Request{
		QType: queryType,
		QName: queryName,
	}
}

func TestRequestHandling(t *testing.T) {
	_, _, lookup := buildLookup(".example.org")

	response, err := lookup.Lookup(buildRequest("example.org", "SOA"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "SOA", response.QType)
	assert.Equal(t, "example.org", response.QName)
	assert.Regexp(t, "dns\\.example\\.org\\. hostmaster\\.example.org\\. \\d+ 1800 3600 7200 5", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("example.org", "NS"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "NS", response.QType)
	assert.Equal(t, "example.org", response.QName)
	assert.Equal(t, "dns.example.org", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("www.example.org", "ANY"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "A", response.QType)
	assert.Equal(t, "www.example.org", response.QName)
	assert.Equal(t, "10.11.12.13", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("www.example.org", "A"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "A", response.QType)
	assert.Equal(t, "www.example.org", response.QName)
	assert.Equal(t, "10.11.12.13", response.Content)
	assert.Equal(t, 5, response.TTL)

	// Allow hostname to be mixed case which is used by Let's Encrypt for a little bit more security
	response, err = lookup.Lookup(buildRequest("wWW.eXaMPlE.oRg", "A"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "A", response.QType)
	assert.Equal(t, "wWW.eXaMPlE.oRg", response.QName)
	assert.Equal(t, "10.11.12.13", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("notexisting.example.org", "A"))
	assert.NotNil(t, err)
	assert.Nil(t, response)

	// Correct Handling of IPv4/IPv6 and ANY/A/AAAA
	response, err = lookup.Lookup(buildRequest("v4.example.org", "ANY"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "A", response.QType)
	assert.Equal(t, "v4.example.org", response.QName)
	assert.Equal(t, "10.10.10.10", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("v4.example.org", "A"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "A", response.QType)
	assert.Equal(t, "v4.example.org", response.QName)
	assert.Equal(t, "10.10.10.10", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("v4.example.org", "AAAA"))
	assert.NotNil(t, err)
	assert.Nil(t, response)

	response, err = lookup.Lookup(buildRequest("v6.example.org", "ANY"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "AAAA", response.QType)
	assert.Equal(t, "v6.example.org", response.QName)
	assert.Equal(t, "2001:db8:85a3::8a2e:370:7334", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("v6.example.org", "AAAA"))
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "AAAA", response.QType)
	assert.Equal(t, "v6.example.org", response.QName)
	assert.Equal(t, "2001:db8:85a3::8a2e:370:7334", response.Content)
	assert.Equal(t, 5, response.TTL)

	response, err = lookup.Lookup(buildRequest("v6.example.org", "A"))
	assert.NotNil(t, err)
	assert.Nil(t, response)
}
