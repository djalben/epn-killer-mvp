package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type PostgresRepo struct {
	Client *sqlx.DB
	Logger *slog.Logger
}

func NewPostgresRepo(client *sqlx.DB, logger *slog.Logger) PostgresRepo {
	return PostgresRepo{
		Client: client,
		Logger: logger,
	}
}
