// Package dbconfig resolves a registry.Database into the dial-ready shape
// every datastore package in go-db needs: host, port, credentials, the
// target database/index name, and a TLS config. It is the one place that
// knows how to turn registry configuration into something you can dial.
package dbconfig

import (
	"crypto/tls"
	"fmt"

	"github.com/ooqls/go-registry"
)

// DialConfig is what's needed to dial a datastore, independent of which one.
type DialConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	TLS      *tls.Config
}

// Resolve turns a registry.Database into a DialConfig. It is pure: it takes
// the config as a value and makes no calls to registry.Get(), so callers
// control where the Database comes from and tests need no global state.
func Resolve(db *registry.Database) (DialConfig, error) {
	if db == nil {
		return DialConfig{}, fmt.Errorf("dbconfig: no database configured")
	}

	tlsCfg, err := db.TLSConfig()
	if err != nil {
		return DialConfig{}, fmt.Errorf("dbconfig: failed to build tls config: %w", err)
	}

	return DialConfig{
		Host:     db.Host,
		Port:     db.Port,
		Username: db.Auth.Username,
		Password: db.Auth.Password,
		Database: db.Database,
		TLS:      tlsCfg,
	}, nil
}

// Addr formats Host and Port as a "host:port" dial address.
func (c DialConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
