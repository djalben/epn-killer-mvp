package app

import (
	"context"

	"github.com/djalben/epn-killer-mvp/internal/application/card"
	"github.com/djalben/epn-killer-mvp/internal/application/commission"
	"github.com/djalben/epn-killer-mvp/internal/application/ticket"
	"github.com/djalben/epn-killer-mvp/internal/application/transaction"
	"github.com/djalben/epn-killer-mvp/internal/application/wallet"
	"github.com/djalben/epn-killer-mvp/internal/config"
	"github.com/djalben/epn-killer-mvp/internal/infrastructure/persistence/postgres"
	"github.com/djalben/epn-killer-mvp/internal/ports"
	"github.com/jmoiron/sqlx"
	"gitlab.com/libs-artifex/wrapper/v2"
)

// Container — главный DI-контейнер проекта.
type Container struct {
	DB *sqlx.DB // храним для корректного закрытия

	// Репозитории
	WalletRepo      ports.WalletRepository
	CardRepo        ports.CardRepository
	TransactionRepo ports.TransactionRepository
	TicketRepo      ports.TicketRepository
	UserRepo        ports.UserRepository
	CommissionRepo  ports.CommissionConfigRepository

	// UseCases
	WalletUseCase      *wallet.UseCase
	CardUseCase        *card.UseCase
	TransactionUseCase *transaction.UseCase
	TicketUseCase      *ticket.UseCase
	CommissionUseCase  *commission.UseCase
}

// NewContainer — создаёт и собирает всё приложение.
func NewContainer(cfg *config.Config) (*Container, error) {
	ctx := context.Background()

	// Подключаемся к БД
	db, err := postgres.Connect(ctx, cfg.PostgresDSN)
	if err != nil {
		return nil, wrapper.Wrap(err)
	}

	// Создаём репозитории (теперь с store)
	walletRepo := postgres.NewWalletRepository(db)
	cardRepo := postgres.NewCardRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	ticketRepo := postgres.NewTicketRepository(db)
	userRepo := postgres.NewUserRepository(db)
	commissionRepo := postgres.NewCommissionConfigRepository(db)

	return &Container{
		DB: db,

		WalletRepo:      walletRepo,
		CardRepo:        cardRepo,
		TransactionRepo: transactionRepo,
		TicketRepo:      ticketRepo,
		UserRepo:        userRepo,
		CommissionRepo:  commissionRepo,

		WalletUseCase:      wallet.NewUseCase(walletRepo, transactionRepo),
		CardUseCase:        card.NewUseCase(cardRepo, walletRepo, transactionRepo),
		TransactionUseCase: transaction.NewUseCase(transactionRepo),
		TicketUseCase:      ticket.NewUseCase(ticketRepo),
		CommissionUseCase:  commission.NewUseCase(commissionRepo),
	}, nil
}

// Close корректно закрывает все соединения с БД.
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}
