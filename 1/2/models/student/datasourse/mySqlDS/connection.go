package mySqlDS

import (
	"database/sql"
	"errors"
	"time"
)

var errEmptyDSN = errors.New("dsn is empty")

func Open(cfg Config) (*sql.DB, error) {
	if cfg.DSN == "" {
		return nil, errEmptyDSN
	}
	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSpan) * time.Second)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
