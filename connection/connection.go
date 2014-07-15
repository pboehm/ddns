package connection

import (
	"crypto/sha1"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
	"time"
)

const HostExpirationSeconds int = 10 * 24 * 60 * 60 // 10 Days

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func OpenConnection() *RedisConnection {
	conn, conn_err := redis.Dial("tcp", ":6379")
	if conn_err != nil {
		log.Fatal(conn_err)
	}

	c := &RedisConnection{conn}
	go c.periodicKeepAlive()

	return c
}

type RedisConnection struct {
	redis.Conn
}

// Every minute send a Ping signal to redis so that we dont run into a timeout
func (self *RedisConnection) periodicKeepAlive() {
	c := time.Tick(1 * time.Minute)

	for _ = range c {
		_, err := self.Do("PING")
		HandleErr(err)
	}
}

func (self *RedisConnection) GetHost(name string) *Host {
	host := Host{Hostname: name}

	if self.HostExist(name) {
		data, err := redis.Values(self.Do("HGETALL", host.Hostname))
		HandleErr(err)

		HandleErr(redis.ScanStruct(data, &host))
	}

	return &host
}

func (self *RedisConnection) SaveHost(host *Host) {
	_, err := self.Do("HMSET", redis.Args{}.Add(host.Hostname).AddFlat(host)...)
	HandleErr(err)

	_, err = self.Do("EXPIRE", host.Hostname, HostExpirationSeconds)
	HandleErr(err)
}

func (self *RedisConnection) HostExist(name string) bool {
	exists, err := redis.Bool(self.Do("EXISTS", name))
	HandleErr(err)

	return exists
}

type Host struct {
	Hostname string `redis:"-"`
	Ip       string `redis:"ip"`
	Token    string `redis:"token"`
}

func (self *Host) GenerateAndSetToken() {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	hash.Write([]byte(self.Hostname))

	self.Token = fmt.Sprintf("%x", hash.Sum(nil))
}

// Returns true when this host has a IPv4 Address and false if IPv6
func (self *Host) IsIPv4() bool {
	if strings.Contains(self.Ip, ".") {
		return true
	}

	return false
}
