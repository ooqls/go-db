package valkey

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/ooqls/go-registry"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
)

var valkeyCfg *registry.Database
var m sync.Mutex = sync.Mutex{}
var c valkey.Client

func Init(db *registry.Database) error {
	m.Lock()
	defer m.Unlock()

	return initValkey(db)
}

func InitDefault() error {
	m.Lock()
	defer m.Unlock()

	if c != nil {
		return nil
	}

	reg := registry.Get()
	if reg.Valkey == nil {
		return fmt.Errorf("no valkey server found in registry")
	}

	return initValkey(reg.Valkey)
}

func initValkey(db *registry.Database) error {
	if c != nil {
		return nil
	}

	valkeyCfg = db

	var tlsCfg *tls.Config
	if db.TLS != nil {
		var err error
		tlsCfg, err = db.TLS.TLSConfig()
		if err != nil {
			return err
		}
	}

	var err error

	c, err = valkey.NewClient(valkey.ClientOption{
		TLSConfig:   tlsCfg,
		Username:    db.Auth.Username,
		Password:    db.Auth.Password,
		InitAddress: []string{fmt.Sprintf("%s:%d", db.Host, db.Port)},
	})
	return err
}

func GetConnection(ctx context.Context) valkey.Client {
	if c == nil {
		err := InitDefault()
		if err != nil {
			panic(fmt.Errorf("failed to init valkey: %v", err))
		}
	} else {
		if err := c.Do(ctx, c.B().Ping().Build()).Error(); err != nil {
			c = nil
			for err != nil {
				m.Lock()
				err := initValkey(valkeyCfg)
				m.Unlock()
				
				if err != nil {
					zap.L().Error("valkey client was disconnected, reconnecting...")
					time.Sleep(time.Second * 3)
				}
			}
		}
	}

	return c
}
