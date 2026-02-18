package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type PostgresRepo struct {
	client *sqlx.DB
	logger *slog.Logger
}

func NewPostgresRepo(client *sqlx.DB, logger *slog.Logger) PostgresRepo {
	return PostgresRepo{
		client: client,
		logger: logger,
	}
}
