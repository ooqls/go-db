package pgx

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/ooqls/go-db/dbconfig"
	"github.com/ooqls/go-db/postgres"
	"github.com/ooqls/go-log"
	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
)

var l *zap.Logger = log.NewLogger("pgx")
var db *pgx.Conn
var m sync.Mutex = sync.Mutex{}

func GetPGX() *pgx.Conn {
	m.Lock()
	defer m.Unlock()

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
			db.Close(context.Background())
		}

		reg := registry.Get()
		cfg, resolveErr := dbconfig.Resolve(reg.Postgres)
		if resolveErr != nil {
			panic(resolveErr)
		}

		db, err = connectPgx(context.Background(), cfg)
		if err != nil {
			panic(err)
		}
	}

	l.Info("PGX database connection established")
	return db
}

func connectPgx(ctx context.Context, cfg dbconfig.DialConfig) (*pgx.Conn, error) {
	connCfg, err := postgres.ConnConfig(cfg)
	if err != nil {
		return nil, err
	}

	conn, err := pgx.ConnectConfig(ctx, connCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	return conn, nil
}

func Init(ctx context.Context, dbCfg *registry.Database) (*pgx.Conn, error) {
	m.Lock()
	defer m.Unlock()

	cfg, err := dbconfig.Resolve(dbCfg)
	if err != nil {
		return nil, err
	}

	conn, err := connectPgx(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func InitDefault() error {
	m.Lock()
	defer m.Unlock()

	reg := registry.Get()
	cfg, err := dbconfig.Resolve(reg.Postgres)
	if err != nil {
		l.Error("failed to resolve postgres config", zap.Error(err))
		return err
	}

	err = initPGX(context.Background(), cfg)
	if err != nil {
		l.Error("failed to initialize default options", zap.Error(err))
		return err
	}
	l.Info("default options initialized successfully")
	return nil
}

func initPGX(ctx context.Context, cfg dbconfig.DialConfig) error {
	var err error
	if db != nil {
		err = postgres.Retry(func() error {
			return db.Ping(ctx)
		})
		if err != nil {
			return err
		}
	}

	if db == nil {
		db, err = connectPgx(ctx, cfg)
		if err != nil {
			return err
		}
	}

	return nil
}
