package shared

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisBackend struct {
	expirationSeconds int
	pool              *redis.Pool
}

func NewRedisBackend(config *Config) *RedisBackend {
	return &RedisBackend{
		expirationSeconds: config.HostExpirationDays * 24 * 60 * 60,
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,

			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", config.RedisHost)
				if err != nil {
					return nil, err
				}
				return c, err
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func (r *RedisBackend) Close() {
	r.pool.Close()
}

func (r *RedisBackend) GetHost(name string) (*Host, error) {
	conn := r.pool.Get()
	defer conn.Close()

	host := Host{Hostname: name}

	var err error
	var data []interface{}

	if data, err = redis.Values(conn.Do("HGETALL", host.Hostname)); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("Host does not exist")
	}

	if err = redis.ScanStruct(data, &host); err != nil {
		return nil, err
	}

	return &host, nil
}

func (r *RedisBackend) SetHost(host *Host) error {
	conn := r.pool.Get()
	defer conn.Close()

	var err error

	if _, err = conn.Do("HMSET", redis.Args{}.Add(host.Hostname).AddFlat(host)...); err != nil {
		return err
	}

	if _, err = conn.Do("EXPIRE", host.Hostname, r.expirationSeconds); err != nil {
		return err
	}

	return nil
}
