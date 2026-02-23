package user

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gitlab.com/libs-artifex/wrapper/v2"

	"github.com/djalben/epn-killer-mvp/internal/entity"
)

// CreateUser создаёт нового пользователя
func (r *Repository) CreateUser(ctx context.Context, email, passwordHash string) (*entity.User, error) {
	const query = `
		INSERT INTO users (email, password_hash, balance, status) 
		VALUES ($1, $2, 0.0000, 'ACTIVE') 
		RETURNING id, email, balance, status, created_at
	`

	var user entity.User
	err := r.Client.QueryRowContext(ctx, query, email, passwordHash).
		Scan(&user.ID, &user.Email, &user.Balance, &user.Status, &user.CreatedAt)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to create user",
			slog.String("email", email),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "user created successfully",
		slog.String("user_id", user.ID),
		slog.String("email", user.Email))

	return &user, nil
}

// GetUserByEmail находит пользователя по email
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	const query = `
		SELECT id, email, password_hash, balance, status, telegram_chat_id, created_at 
		FROM users 
		WHERE email = $1
	`

	var user entity.User
	err := r.Client.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrapper.Wrap(errors.New("user not found"))
		}
		r.Logger.ErrorContext(ctx, "failed to get user by email",
			slog.String("email", email),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	return &user, nil
}

// GetUserByID находит пользователя по ID
func (r *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	const query = `
		SELECT id, email, password_hash, balance, status, telegram_chat_id, created_at 
		FROM users 
		WHERE id = $1
	`

	var user entity.User
	err := r.Client.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrapper.Wrap(errors.New("user not found"))
		}
		r.Logger.ErrorContext(ctx, "failed to get user by id",
			slog.String("user_id", id.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	return &user, nil
}

// UpdateTelegramChatID обновляет Telegram Chat ID пользователя
func (r *Repository) UpdateTelegramChatID(ctx context.Context, userID uuid.UUID, chatID int64) error {
	const query = `
		UPDATE users 
		SET telegram_chat_id = $1 
		WHERE id = $2
	`

	_, err := r.Client.ExecContext(ctx, query, chatID, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to update telegram chat id",
			slog.String("user_id", userID.String()),
			slog.Int64("chat_id", chatID),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "telegram chat id updated",
		slog.String("user_id", userID.String()),
		slog.Int64("chat_id", chatID))

	return nil
}
