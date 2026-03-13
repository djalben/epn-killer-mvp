package ticket

import (
	"context"

	"github.com/djalben/epn-killer-mvp/internal/domain"
	"github.com/djalben/epn-killer-mvp/internal/ports"
)

type UseCase struct {
	ticketRepo ports.TicketRepository
}

func NewUseCase(tr ports.TicketRepository) *UseCase {
	return &UseCase{ticketRepo: tr}
}

func (uc *UseCase) Create(ctx context.Context, userID domain.UUID, subject, message string, tgChatID *int64) (*domain.Ticket, error) {
	t := domain.NewTicket(userID, subject, message, tgChatID)
	return t, uc.ticketRepo.Save(ctx, t)
}

func (uc *UseCase) Take(ctx context.Context, ticketID domain.UUID, adminID domain.UUID) error {
	t, err := uc.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}
	t.Take(adminID)
	return uc.ticketRepo.Update(ctx, t)
}

func (uc *UseCase) Close(ctx context.Context, ticketID domain.UUID, reply string) error {
	t, err := uc.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}
	t.Close(reply)
	return uc.ticketRepo.Update(ctx, t)
}
