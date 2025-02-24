package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx"
	"github.com/ooqls/go-db/postgres"
	"github.com/ooqls/go-log"
	"go.uber.org/zap"
)

var l *zap.Logger = log.NewLogger("pgx")
var db *pgx.Conn

func GetDBX() *pgx.Conn {
	var err error
	if db != nil {
		err = postgres.Retry(func() error {
			return db.Ping(context.Background())
		})
		if err != nil {
			l.Error("failed to ping connection", zap.Error(err))
		}
	}

	if db == nil || err != nil {
		if db != nil {
			db.Close()
		}

		db, err = connectPgx(postgres.GetRegistryOptions())
		if err != nil {
			panic(err)
		}
	}

	l.Info("PGX database connection established")
	return db
}

func connectPgx(opt postgres.Options) (*pgx.Conn, error) {

	conn, err := pgx.Connect(pgx.ConnConfig{
		Host:      opt.Host,
		Port:      uint16(opt.Port),
		User:      opt.User,
		Password:  opt.Pw,
		Database:  opt.DB,
		TLSConfig: opt.Tls,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	return conn, nil
}

func Init(opt postgres.Options) (*pgx.Conn, error) {
	conn, err := connectPgx(opt)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func InitDefault() error {
	opts := postgres.GetRegistryOptions()
	_, err := Init(opts)
	if err != nil {

		l.Error("failed to initialize default options", zap.Error(err))
		return err
	}
	l.Info("default options initialized successfully")
	return nil
}
