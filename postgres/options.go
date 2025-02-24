package postgres

import (
	"crypto/tls"

	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
)

type Options struct {
	Host string
	Port int
	User string
	DB   string
	Pw   string
	Tls  *tls.Config
}

var dbName string = "postgres"

func GetRegistryOptions() Options {
	reg := registry.Get()
	var tlsCfg *tls.Config
	var err error
	if reg.Postgres.TLS != nil {
		tlsCfg, err = reg.Postgres.TLS.TLSConfig()
		if err != nil {
			l.Error("failed to get TLS config", zap.Error(err))
			panic(err)
		}
	}

	return Options{
		Host: reg.Postgres.Host,
		Port: reg.Postgres.Port,
		User: reg.Postgres.Auth.Username,
		Pw:   reg.Postgres.Auth.Password,
		DB:   dbName,
		Tls:  tlsCfg,
	}
}
