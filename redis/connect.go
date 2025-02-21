package redis

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"sync"

	"github.com/ooqls/go-registry"
	"github.com/redis/go-redis/v9"
)

var pool *redis.Client
var m sync.Mutex = sync.Mutex{}

func initDefault() error {
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

	redisDb, err := strconv.Atoi(reg.Redis.Database)
	if err != nil {
		return fmt.Errorf("failed to convert database string to int: %v", err)
	}

	pool = redis.NewClient(&redis.Options{

		Addr:      fmt.Sprintf("%s:%d", reg.Redis.Host, reg.Redis.Port),
		Password:  reg.Redis.Auth.Password,
		Username:  reg.Redis.Auth.Username,
		DB:        redisDb,
		TLSConfig: tlsCfg,
		PoolSize:  3,
	})

	return nil
}

func GetConnection() *redis.Client {
	m.Lock()
	defer m.Unlock()

	if pool == nil {
		if err := initDefault(); err != nil {
			panic(err)
		}

	}
	return pool
}
