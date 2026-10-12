package sqlx

import (
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/ooqls/go-db/dbconfig"
	"github.com/ooqls/go-db/postgres"
	log "github.com/ooqls/go-log"
	"github.com/ooqls/go-registry"
	"go.uber.org/zap"
)

var db *sqlx.DB
var l *zap.Logger = log.NewLogger("db")
var m sync.Mutex = sync.Mutex{}

func GetSQLX() *sqlx.DB {
	var err error
	if db != nil {
		err = postgres.Retry(func() error {
			return db.Ping()
		})
		if err != nil {
			l.Error("failed to ensure connection", zap.Error(err))
		}
	}

	if db == nil || err != nil {
		if db != nil {
			db.Close()
		}

		reg := registry.Get()
		cfg, resolveErr := dbconfig.Resolve(reg.Postgres)
		if resolveErr != nil {
			panic(resolveErr)
		}

		db, err = connectSqlx(cfg)
		if err != nil {
			panic(err)
		}
	}

	l.Info("SQLX database connection established")
	return db
}

func connectSqlx(cfg dbconfig.DialConfig) (*sqlx.DB, error) {
	connCfg, err := postgres.ConnConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build postgres connection config: %v", err)
	}

	sqlDB := stdlib.OpenDB(*connCfg)
	dbCon := sqlx.NewDb(sqlDB, "pgx")

	db = dbCon

	return dbCon, nil
}

func Init(dbCfg *registry.Database) (*sqlx.DB, error) {
	m.Lock()
	defer m.Unlock()

	cfg, err := dbconfig.Resolve(dbCfg)
	if err != nil {
		return nil, err
	}

	err = initSQLX(cfg)
	return db, err
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

	err = initSQLX(cfg)
	if err != nil {
		l.Error("failed to initialize default options", zap.Error(err))
		return err
	}
	l.Info("default options initialized successfully")
	return nil
}

func initSQLX(cfg dbconfig.DialConfig) error {
	var err error
	if db != nil {
		err = postgres.Retry(func() error {
			return db.Ping()
		})
		if err != nil {
			l.Error("failed to ensure connection", zap.Error(err))
		}
	}

	if db == nil || err != nil {
		if db != nil {
			db.Close()
		}

		_, err = connectSqlx(cfg)
		if err != nil {
			return err
		}
	}

	l.Info("SQLX database connection established")
	return nil
}
