package valkey

import (
	"sync"

	"github.com/ooqls/go-registry"
	"github.com/valkey-io/valkey-go"
)

var valkeyCfg *registry.Database
var m sync.Mutex = sync.Mutex{}
var cFunc func(db *registry.Database) (*valkey.Client, error)
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



	return initValkey(valkeyCfg)
}

func initValkey(db *registry.Database) error {
	m.Lock()
	defer m.Unlock()

	if c != nil {
		return nil
	}

	valkeyCfg = db
	var err error
	tlsCfg, err := db.TLS.TLSConfig()
	if err != nil {
		return err
	}

	c, err = valkey.NewClient(valkey.ClientOption{
		TLSConfig: tlsCfg,
		Username: db.Auth.Username,
		Password: db.Auth.Password,
		InitAddress: []string{db.Host},
	})
	return err
}
