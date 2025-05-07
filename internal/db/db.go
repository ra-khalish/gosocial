package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New(addr string, maxOpenConn, maxIdleConn int, maxIdleTime string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.New(ctx, addr)
	// db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	// db.SetMaxOpenConns(maxOpenConn)
	// db.SetMaxIdleConns(maxIdleConn)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.Config().MaxConnIdleTime = duration
	db.Config().MaxConns = int32(maxOpenConn)
	// db.SetConnMaxIdleTime(duration)

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	if err = db.Ping(ctx); err != nil {
		return nil, err
	}
	return db, nil

}
