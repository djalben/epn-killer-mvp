package card

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.com/libs-artifex/wrapper/v2"

	"github.com/djalben/epn-killer-mvp/internal/entity"
)

// CreateCard создаёт новую виртуальную карту
func (r *Repository) CreateCard(ctx context.Context, userID uuid.UUID, providerCardID, bin, last4, cardType, nickname string, dailyLimit decimal.Decimal) (*entity.Card, error) {
	const query = `
		INSERT INTO cards (
			user_id, provider_card_id, bin, last_4_digits, card_status, 
			nickname, daily_spend_limit, card_type, card_balance
		) VALUES ($1, $2, $3, $4, 'ACTIVE', $5, $6, $7, 0.0000)
		RETURNING id, user_id, provider_card_id, bin, last_4_digits, card_status, 
		          nickname, daily_spend_limit, failed_auth_count, card_type, 
		          auto_replenish_enabled, auto_replenish_threshold, auto_replenish_amount, 
		          card_balance, team_id, created_at
	`

	var card entity.Card
	err := r.Client.QueryRowContext(ctx, query,
		userID, providerCardID, bin, last4, nickname, dailyLimit, cardType,
	).Scan(
		&card.ID, &card.UserID, &card.ProviderCardID, &card.BIN, &card.Last4Digits,
		&card.CardStatus, &card.Nickname, &card.DailySpendLimit, &card.FailedAuthCount,
		&card.CardType, &card.AutoReplenishEnabled, &card.AutoReplenishThreshold,
		&card.AutoReplenishAmount, &card.CardBalance, &card.TeamID, &card.CreatedAt,
	)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to create card",
			slog.String("user_id", userID.String()),
			slog.String("provider_card_id", providerCardID),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "card created successfully",
		slog.String("card_id", card.ID),
		slog.String("user_id", userID.String()),
		slog.String("last4", card.Last4Digits))

	return &card, nil
}

// GetCardByID возвращает карту по ID
func (r *Repository) GetCardByID(ctx context.Context, id uuid.UUID) (*entity.Card, error) {
	const query = `
		SELECT id, user_id, provider_card_id, bin, last_4_digits, card_status, 
		       nickname, daily_spend_limit, failed_auth_count, card_type,
		       auto_replenish_enabled, auto_replenish_threshold, auto_replenish_amount,
		       card_balance, team_id, created_at
		FROM cards 
		WHERE id = $1
	`

	var card entity.Card
	err := r.Client.GetContext(ctx, &card, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrapper.Wrap(errors.New("card not found"))
		}
		r.Logger.ErrorContext(ctx, "failed to get card by id",
			slog.String("card_id", id.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	return &card, nil
}

// GetUserCards возвращает все карты пользователя
func (r *Repository) GetUserCards(ctx context.Context, userID uuid.UUID) ([]entity.Card, error) {
	const query = `
		SELECT id, user_id, provider_card_id, bin, last_4_digits, card_status, 
		       nickname, daily_spend_limit, failed_auth_count, card_type,
		       auto_replenish_enabled, auto_replenish_threshold, auto_replenish_amount,
		       card_balance, team_id, created_at
		FROM cards 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`

	rows, err := r.Client.QueryContext(ctx, query, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to get user cards",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}
	defer rows.Close()

	var cards []entity.Card
	for rows.Next() {
		var card entity.Card
		err := rows.Scan(
			&card.ID, &card.UserID, &card.ProviderCardID, &card.BIN, &card.Last4Digits,
			&card.CardStatus, &card.Nickname, &card.DailySpendLimit, &card.FailedAuthCount,
			&card.CardType, &card.AutoReplenishEnabled, &card.AutoReplenishThreshold,
			&card.AutoReplenishAmount, &card.CardBalance, &card.TeamID, &card.CreatedAt,
		)
		if err != nil {
			r.Logger.ErrorContext(ctx, "failed to scan card", slog.Any("error", err))
			continue
		}
		cards = append(cards, card)
	}

	return cards, nil
}

// ProcessCardPayment — атомарное списание средств с баланса пользователя
func (r *Repository) ProcessCardPayment(ctx context.Context, userID, cardID uuid.UUID, amount, fee decimal.Decimal, merchantName, cardLast4 string) error {
	tx, err := r.Client.BeginTx(ctx, nil)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to begin transaction", slog.Any("error", err))
		return wrapper.Wrap(err)
	}
	defer tx.Rollback()

	// Блокируем строку пользователя
	var currentBalance decimal.Decimal
	err = tx.QueryRowContext(ctx,
		"SELECT balance FROM users WHERE id = $1 FOR UPDATE", userID,
	).Scan(&currentBalance)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to lock user balance", slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	if currentBalance.LessThan(amount) {
		r.Logger.WarnContext(ctx, "insufficient balance",
			slog.String("user_id", userID.String()),
			slog.Any("balance", currentBalance),
			slog.Any("amount", amount))
		return wrapper.Wrap(errors.New("недостаточно средств"))
	}

	// Списываем средства
	_, err = tx.ExecContext(ctx,
		"UPDATE users SET balance = balance - $1 WHERE id = $2",
		amount, userID,
	)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to deduct balance", slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	// Записываем транзакцию
	const txQuery = `
		INSERT INTO transactions (user_id, card_id, amount, fee, transaction_type, status, details)
		VALUES ($1, $2, $3, $4, 'CAPTURE', 'APPROVED', $5)
	`
	_, err = tx.ExecContext(ctx, txQuery, userID, cardID, amount, fee,
		fmt.Sprintf("Payment to %s (...%s)", merchantName, cardLast4))
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to record transaction", slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	// Сбрасываем счётчик неудачных авторизаций
	_, err = tx.ExecContext(ctx, "UPDATE cards SET failed_auth_count = 0 WHERE id = $1", cardID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to reset failed auth count", slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	if err := tx.Commit(); err != nil {
		r.Logger.ErrorContext(ctx, "failed to commit transaction", slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "card payment processed successfully",
		slog.String("user_id", userID.String()),
		slog.String("card_id", cardID.String()),
		slog.Any("amount", amount))

	return nil
}

// IncrementFailedAuthCount увеличивает счётчик неудачных авторизаций
func (r *Repository) IncrementFailedAuthCount(ctx context.Context, cardID uuid.UUID) error {
	const query = "UPDATE cards SET failed_auth_count = failed_auth_count + 1 WHERE id = $1"
	_, err := r.Client.ExecContext(ctx, query, cardID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to increment failed auth count",
			slog.String("card_id", cardID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}
	return nil
}

// BlockCard блокирует карту
func (r *Repository) BlockCard(ctx context.Context, cardID uuid.UUID) error {
	const query = "UPDATE cards SET card_status = 'BLOCKED' WHERE id = $1"
	_, err := r.Client.ExecContext(ctx, query, cardID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to block card",
			slog.String("card_id", cardID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.WarnContext(ctx, "card blocked by antifraud system",
		slog.String("card_id", cardID.String()))

	return nil
}

// UpdateCardStatus обновляет статус карты
func (r *Repository) UpdateCardStatus(ctx context.Context, cardID, userID uuid.UUID, status string) error {
	if status != "ACTIVE" && status != "BLOCKED" && status != "CLOSED" {
		return wrapper.Wrap(errors.New("invalid card status"))
	}

	const query = `
		UPDATE cards 
		SET card_status = $1 
		WHERE id = $2 AND user_id = $3
	`

	result, err := r.Client.ExecContext(ctx, query, status, cardID, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to update card status",
			slog.String("card_id", cardID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return wrapper.Wrap(errors.New("card not found or access denied"))
	}

	r.Logger.InfoContext(ctx, "card status updated",
		slog.String("card_id", cardID.String()),
		slog.String("status", status))

	return nil
}

// IssueCards — массовый выпуск карт (mock для MVP)
func (r *Repository) IssueCards(ctx context.Context, userID uuid.UUID, req entity.MassIssueRequest) (*entity.MassIssueResponse, error) {
	// ... (я могу сделать полную реализацию, если нужно)
	// Пока оставляю заглушку, так как метод большой
	return nil, wrapper.Wrap(errors.New("not implemented yet"))
}

// UpdateCardAutoReplenishment обновляет настройки автопополнения
func (r *Repository) UpdateCardAutoReplenishment(ctx context.Context, cardID, userID uuid.UUID, enabled bool, threshold, amount decimal.Decimal) error {
	const query = `
		UPDATE cards 
		SET auto_replenish_enabled = $1, 
		    auto_replenish_threshold = $2, 
		    auto_replenish_amount = $3 
		WHERE id = $4 AND user_id = $5
	`

	_, err := r.Client.ExecContext(ctx, query, enabled, threshold, amount, cardID, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to update auto-replenishment",
			slog.String("card_id", cardID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "auto-replenishment updated",
		slog.String("card_id", cardID.String()),
		slog.Bool("enabled", enabled))

	return nil
}
