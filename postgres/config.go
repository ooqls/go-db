package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/ooqls/go-db/dbconfig"
)

// ConnConfig turns a dbconfig.DialConfig into a pgx connection config ready
// to dial Postgres, with TLS applied exactly as dbconfig resolved it.
func ConnConfig(cfg dbconfig.DialConfig) (*pgx.ConnConfig, error) {
	sslmode := "disable"
	if cfg.TLS != nil {
		sslmode = "require"
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, sslmode)

	connCfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection config: %w", err)
	}

	// ParseConfig derives a default TLSConfig from sslmode; override it with
	// the one dbconfig resolved so cert/CA/InsecureSkipVerify settings apply.
	if cfg.TLS != nil {
		connCfg.TLSConfig = cfg.TLS
	}

	return connCfg, nil
}
