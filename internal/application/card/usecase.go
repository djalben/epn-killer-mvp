package card

import (
	"context"

	"github.com/djalben/epn-killer-mvp/internal/domain"
	"github.com/djalben/epn-killer-mvp/internal/ports"
)

type UseCase struct {
	cardRepo   ports.CardRepository
	walletRepo ports.WalletRepository
	txRepo     ports.TransactionRepository
}

func NewUseCase(cr ports.CardRepository, wr ports.WalletRepository, tr ports.TransactionRepository) *UseCase {
	return &UseCase{cardRepo: cr, walletRepo: wr, txRepo: tr}
}

func (uc *UseCase) BuyCard(ctx context.Context, userID domain.UUID, cardType domain.CardType) (*domain.Card, error) {
	card, err := domain.NewCard(userID, cardType, "TEMP_PROVIDER_ID")
	if err != nil {
		return nil, err
	}

	if err := uc.cardRepo.Save(ctx, card); err != nil {
		return nil, err
	}

	tx := domain.NewTransaction(userID, &card.ID, domain.NewNumeric(2.00), domain.NewNumeric(0), "CARD_ISSUE", "COMPLETED", "Выпуск карты")
	return card, uc.txRepo.Save(ctx, tx)
}

func (uc *UseCase) TopUpCard(ctx context.Context, userID domain.UUID, cardID domain.UUID, amount domain.Numeric) error {
	wallet, err := uc.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	card, err := uc.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return err
	}

	if err := wallet.Withdraw(amount); err != nil {
		return err
	}
	card.Balance = card.Balance.Add(amount)

	if err := uc.walletRepo.Update(ctx, wallet); err != nil {
		return err
	}
	if err := uc.cardRepo.Update(ctx, card); err != nil {
		return err
	}

	tx := domain.NewTransaction(userID, &cardID, amount, domain.NewNumeric(0), "TOPUP_CARD", "COMPLETED", "Пополнение карты")
	return uc.txRepo.Save(ctx, tx)
}

func (uc *UseCase) ToggleAutoTopUp(ctx context.Context, cardID domain.UUID, enabled bool, below, amount domain.Numeric) error {
	card, err := uc.cardRepo.GetByID(ctx, cardID)
	if err != nil {
		return err
	}

	card.AutoTopUpEnabled = enabled
	card.AutoTopUpBelow = below
	card.AutoTopUpAmount = amount

	return uc.cardRepo.Update(ctx, card)
}
