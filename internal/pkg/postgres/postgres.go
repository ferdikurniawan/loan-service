package postgres

import (
	"context"
	"database/sql"
	"time"
)

type Postgres struct {
	DB *sql.DB
}

type Config struct {
	Dsn           string
	MaxConn       int
	MaxIdle       int
	DataDogTracer bool //trace to be on or off
}

const (
	// postgres driver name
	postgres = "postgres"
)

func New(cfg *Config) (*Postgres, error) {
	var db *sql.DB
	var err error
	db, err = sql.Open(postgres, cfg.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxConn)
	db.SetMaxIdleConns(cfg.MaxIdle)

	//ping db to ensure the connection is alive and working
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return &Postgres{db}, nil
}

func (p *Postgres) Close() {
	p.DB.Close()
}
