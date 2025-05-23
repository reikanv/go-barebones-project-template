package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	Pool *pgxpool.Pool
}

func NewPGX(ctx context.Context, dsn string) (*Db, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return &Db{Pool: pool}, nil
}

func (db *Db) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}
