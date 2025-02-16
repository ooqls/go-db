package redis

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/braumsmilk/go-registry"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool
var m sync.Mutex = sync.Mutex{}

func InitDefault() error {
	m.Lock()
	defer m.Unlock()

	if pool != nil {
		return nil
	}

	reg := registry.Get()
	if reg.Redis == nil {
		return fmt.Errorf("no redis server found in registry")
	}
	var tlsCfg *tls.Config
	if reg.Redis.TLS != nil {
		var err error
		tlsCfg, err = reg.Redis.TLS.TLSConfig()
		if err != nil {
			return fmt.Errorf("failed to load tls config for redis: %v", err)
		}
	}

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				fmt.Sprintf("%s:%d", reg.Redis.Host, reg.Redis.Port),
				redis.DialUseTLS(reg.Redis.TLS != nil),
				redis.DialPassword(reg.Redis.Auth.Password),
				redis.DialTLSConfig(tlsCfg))
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func GetConnection() redis.Conn {
	m.Lock()
	defer m.Unlock()
	
	return pool.Get()
}
