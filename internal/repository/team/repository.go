package team

import (
	"github.com/djalben/epn-killer-mvp/internal/repository"
)

type Repository struct {
	*repository.PostgresRepo
}

func New(repo *repository.PostgresRepo) *Repository {
	return &Repository{PostgresRepo: repo}
}
