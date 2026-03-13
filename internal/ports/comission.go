package ports

import (
	"context"

	"github.com/djalben/epn-killer-mvp/internal/domain"
)

type CommissionConfigRepository interface {
	GetByKey(ctx context.Context, key string) (*domain.CommissionConfig, error)
	Update(ctx context.Context, config *domain.CommissionConfig) error
	ListAll(ctx context.Context) ([]*domain.CommissionConfig, error)
}
