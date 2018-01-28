package shared

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"
)

type Host struct {
	Hostname string `redis:"-"`
	Ip       string `redis:"ip"`
	Token    string `redis:"token"`
}

func (h *Host) GenerateAndSetToken() {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	hash.Write([]byte(h.Hostname))

	h.Token = fmt.Sprintf("%x", hash.Sum(nil))
}

func (h *Host) IsIPv4() bool {
	if strings.Contains(h.Ip, ".") {
		return true
	}

	return false
}

type HostBackend interface {
	GetHost(string) (*Host, error)

	SetHost(*Host) error
}
