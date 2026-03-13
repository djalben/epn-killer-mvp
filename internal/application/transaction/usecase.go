package transaction

import (
	"context"
	"time"

	"github.com/djalben/epn-killer-mvp/internal/domain"
	"github.com/djalben/epn-killer-mvp/internal/ports"
)

type UseCase struct {
	txRepo ports.TransactionRepository
}

func NewUseCase(tr ports.TransactionRepository) *UseCase {
	return &UseCase{txRepo: tr}
}

// GetWalletTransactions — история по кошельку
func (uc *UseCase) GetWalletTransactions(ctx context.Context, userID domain.UUID, from, to time.Time) ([]*domain.Transaction, error) {
	return uc.txRepo.GetWalletTransactions(ctx, userID, from, to)
}

// GetCardTransactions — история по одной карте
func (uc *UseCase) GetCardTransactions(ctx context.Context, cardID domain.UUID, from, to time.Time) ([]*domain.Transaction, error) {
	return uc.txRepo.GetCardTransactions(ctx, cardID, from, to)
}
