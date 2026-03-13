package commission

import (
	"context"

	"github.com/djalben/epn-killer-mvp/internal/domain"
	"github.com/djalben/epn-killer-mvp/internal/ports"
)

type UseCase struct {
	configRepo ports.CommissionConfigRepository
}

func NewUseCase(cr ports.CommissionConfigRepository) *UseCase {
	return &UseCase{configRepo: cr}
}

func (uc *UseCase) GetByKey(ctx context.Context, key string) (*domain.CommissionConfig, error) {
	return uc.configRepo.GetByKey(ctx, key)
}

func (uc *UseCase) Update(ctx context.Context, cfg *domain.CommissionConfig) error {
	return uc.configRepo.Update(ctx, cfg)
}

func (uc *UseCase) ListAll(ctx context.Context) ([]*domain.CommissionConfig, error) {
	return uc.configRepo.ListAll(ctx)
}
