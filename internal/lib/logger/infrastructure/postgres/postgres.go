package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Connect postgres database.
func Connect(ctx context.Context, dsn string) (*sqlx.DB, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.New DSN = %s error = %w", dsn, err)
	}

	pgxPool := stdlib.OpenDBFromPool(pool)

	database := sqlx.NewDb(pgxPool, "pgx")

	err = database.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("PingContext DSN = %s error = %w", dsn, err)
	}

	return database, nil
}
