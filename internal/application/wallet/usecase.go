package wallet

import (
	"context"

	"github.com/djalben/epn-killer-mvp/internal/domain"
	"github.com/djalben/epn-killer-mvp/internal/ports"
)

type UseCase struct {
	walletRepo ports.WalletRepository
	txRepo     ports.TransactionRepository
}

func NewUseCase(wr ports.WalletRepository, tr ports.TransactionRepository) *UseCase {
	return &UseCase{walletRepo: wr, txRepo: tr}
}

func (uc *UseCase) TopUpWallet(ctx context.Context, userID domain.UUID, amount domain.Numeric) error {
	wallet, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	if err := wallet.TopUp(amount); err != nil {
		return err
	}

	if err := uc.walletRepo.Update(ctx, wallet); err != nil {
		return err
	}

	tx := domain.NewTransaction(userID, nil, amount, domain.NewNumeric(0), "TOPUP_WALLET", "COMPLETED", "Пополнение по СБП")
	return uc.txRepo.Save(ctx, tx)
}

func (uc *UseCase) GetBalance(ctx context.Context, userID domain.UUID) (domain.Numeric, error) {
	w, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return domain.NewNumeric(0), err
	}
	return w.Balance, nil
}
