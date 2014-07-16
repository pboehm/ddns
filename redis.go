package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"strings"
	"time"
)

const HostExpirationSeconds int = 10 * 24 * 60 * 60 // 10 Days

type RedisConnection struct {
	*redis.Pool
}

func OpenConnection(server string) *RedisConnection {
	return &RedisConnection{newPool(server)}
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle: 3,

		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func (self *RedisConnection) GetHost(name string) *Host {
	conn := self.Get()
	defer conn.Close()

	host := Host{Hostname: name}

	if self.HostExist(name) {
		data, err := redis.Values(conn.Do("HGETALL", host.Hostname))
		HandleErr(err)

		HandleErr(redis.ScanStruct(data, &host))
	}

	return &host
}

func (self *RedisConnection) SaveHost(host *Host) {
	conn := self.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", redis.Args{}.Add(host.Hostname).AddFlat(host)...)
	HandleErr(err)

	_, err = conn.Do("EXPIRE", host.Hostname, HostExpirationSeconds)
	HandleErr(err)
}

func (self *RedisConnection) HostExist(name string) bool {
	conn := self.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", name))
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
