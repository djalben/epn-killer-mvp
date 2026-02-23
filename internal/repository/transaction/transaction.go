package transaction

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.com/libs-artifex/wrapper/v2"

	"github.com/djalben/epn-killer-mvp/internal/entity"
	"github.com/djalben/epn-killer-mvp/internal/notification"
)

// ProcessDeposit обрабатывает пополнение баланса пользователя и записывает транзакцию
func (r *Repository) ProcessDeposit(ctx context.Context, userID uuid.UUID, amount decimal.Decimal) (*entity.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, wrapper.Wrap(errors.New("сумма пополнения должна быть положительной"))
	}

	tx, err := r.Client.BeginTx(ctx, nil)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to begin transaction for deposit",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}
	defer tx.Rollback()

	// 1. Увеличиваем баланс пользователя атомарно
	const updateBalanceQuery = `
		UPDATE users 
		SET balance = balance + $1 
		WHERE id = $2
	`
	_, err = tx.ExecContext(ctx, updateBalanceQuery, amount, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to update user balance",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	// 2. Создаём запись транзакции
	const insertTxQuery = `
		INSERT INTO transactions (
			user_id, amount, fee, transaction_type, status, details, executed_at
		) VALUES ($1, $2, $3, 'FUND', 'APPROVED', $4, NOW())
		RETURNING id, user_id, card_id, amount, fee, transaction_type, status, details, executed_at
	`

	var m transactionModel
	err = tx.QueryRowContext(ctx, insertTxQuery,
		userID,
		amount,
		decimal.Zero,
		fmt.Sprintf("Deposit via API. Amount: %s", amount.String()),
	).Scan(
		&m.ID,
		&m.UserID,
		&m.CardID,
		&m.Amount,
		&m.Fee,
		&m.TransactionType,
		&m.Status,
		&m.Details,
		&m.ExecutedAt,
	)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to record deposit transaction",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	// 3. Коммит транзакции
	if err := tx.Commit(); err != nil {
		r.Logger.ErrorContext(ctx, "failed to commit deposit transaction",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	// 4. Конвертируем внутреннюю модель в entity для возврата
	result := m.toEntity()

	r.Logger.InfoContext(ctx, "deposit processed successfully",
		slog.String("user_id", userID.String()),
		slog.Any("amount", amount),
		slog.String("transaction_id", result.TransactionID))

	// 5. Отправка уведомления в Telegram (асинхронно)
	go func() {
		user, err := r.userRepo.GetUserByID(ctx, userID)
		if err == nil && user.TelegramChatID != nil {
			message := fmt.Sprintf("✅ Пополнение успешно!\n\nСумма: $%s\nНовый баланс: $%s",
				amount.String(),
				user.Balance.String())
			notification.SendTelegramMessage(user.TelegramChatID, message)
		}
	}()

	return &result, nil
}
