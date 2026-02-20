package team

import (
	"github.com/djalben/epn-killer-mvp/internal/repository"
	"github.com/djalben/epn-killer-mvp/internal/repository/user"
)

type Repository struct {
	*repository.PostgresRepo
	userRepo user.Repository // ← поле для user-репозитория
}

func New(repo *repository.PostgresRepo, userRepo user.Repository) *Repository {
	return &Repository{
		PostgresRepo: repo,
		userRepo:     userRepo,
	}
}
