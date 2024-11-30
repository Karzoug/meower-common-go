package postgresql

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	*pgxpool.Pool
}

func NewDB(ctx context.Context, cfg Config) (DB, error) {
	const op = "postgresql: new db"

	pgxCfg, err := pgxpool.ParseConfig(cfg.URI)
	if err != nil {
		return DB{}, fmt.Errorf("%s: %w", op, err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return DB{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(ctx); err != nil {
		return DB{}, fmt.Errorf("%s: %w", op, err)
	}

	return DB{Pool: pool}, nil
}

func (db *DB) Close(_ context.Context) error {
	db.Pool.Close()
	return nil
}
