package dbconfig

import (
	"testing"

	"github.com/ooqls/go-registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve_NilDatabase(t *testing.T) {
	_, err := Resolve(nil)
	require.Error(t, err)
}

func TestResolve_FieldPassThrough(t *testing.T) {
	db := &registry.Database{
		Database: "widgets",
		Server: registry.Server{
			Host: "db.internal",
			Port: 5432,
			Auth: registry.Auth{
				Username: "alice",
				Password: "s3cret",
			},
		},
	}

	cfg, err := Resolve(db)
	require.NoError(t, err)
	assert.Equal(t, "db.internal", cfg.Host)
	assert.Equal(t, 5432, cfg.Port)
	assert.Equal(t, "alice", cfg.Username)
	assert.Equal(t, "s3cret", cfg.Password)
	assert.Equal(t, "widgets", cfg.Database)
	assert.Equal(t, "db.internal:5432", cfg.Addr())
}

func TestResolve_NoTLSConfigured(t *testing.T) {
	db := &registry.Database{
		Server: registry.Server{Host: "localhost", Port: 1234},
	}

	cfg, err := Resolve(db)
	require.NoError(t, err)
	assert.Nil(t, cfg.TLS)
}

func TestResolve_TLSDisabled(t *testing.T) {
	db := &registry.Database{
		Server: registry.Server{
			Host: "localhost",
			Port: 1234,
			TLS:  &registry.TLSConfig{Enabled: false},
		},
	}

	cfg, err := Resolve(db)
	require.NoError(t, err)
	assert.Nil(t, cfg.TLS)
}

func TestResolve_TLSEnabled(t *testing.T) {
	db := &registry.Database{
		Server: registry.Server{
			Host: "localhost",
			Port: 1234,
			TLS: &registry.TLSConfig{
				Enabled:               true,
				InsecureSkipTLSVerify: true,
			},
		},
	}

	cfg, err := Resolve(db)
	require.NoError(t, err)
	require.NotNil(t, cfg.TLS)
	assert.True(t, cfg.TLS.InsecureSkipVerify)
}

func TestResolve_TLSError(t *testing.T) {
	db := &registry.Database{
		Server: registry.Server{
			Host: "localhost",
			Port: 1234,
			TLS: &registry.TLSConfig{
				Enabled:  true,
				CertPath: "cert.pem",
				// KeyPath deliberately omitted: cert/key must be set together.
			},
		},
	}

	_, err := Resolve(db)
	require.Error(t, err)
}
