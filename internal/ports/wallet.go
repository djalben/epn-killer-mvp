package ports

import (
	"context"

	"github.com/djalben/epn-killer-mvp/internal/domain"
)

type WalletRepository interface {
	GetByUserID(ctx context.Context, userID domain.UUID) (*domain.Wallet, error)
	Update(ctx context.Context, wallet *domain.Wallet) error
}
