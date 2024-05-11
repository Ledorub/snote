package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func CreatePool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("DSN parsing: %w", err)
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connection pool: %w", err)
	}
	if err = CheckConnection(ctx, pool); err != nil {
		return nil, err
	}
	return pool, nil
}

func CheckConnection(ctx context.Context, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("connection pool: %w", err)
	}
	defer conn.Release()

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := conn.Ping(pingCtx); err != nil {
		return fmt.Errorf("DB ping: %w", err)
	}
	return nil
}

func BuildDSN(host string, port uint64, name, user, password string) string {
	return fmt.Sprintf("postgres://%s:%d@%s:%s/%s", host, port, user, password, name)
}
