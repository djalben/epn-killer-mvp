package apikey

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"gitlab.com/libs-artifex/wrapper/v2"
)

// GenerateAPIKey генерирует новый ключ и сохраняет его для пользователя (UPSERT: старый удаляется)
func (r *Repository) GenerateAPIKey(ctx context.Context, userID uuid.UUID) (string, error) {
	// 1. Генерируем новый ключ
	hexKey, err := generateRandomString(16)
	if err != nil {
		return "", wrapper.Wrap(err)
	}
	apiKeyUUID := formatAsUUID(hexKey)

	// 2. Удаляем старый ключ пользователя (если есть)
	_, err = r.Client.ExecContext(ctx, "DELETE FROM api_keys WHERE user_id = $1", userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to delete old API key",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return "", wrapper.Wrap(err)
	}

	// 3. Вставляем новый ключ
	const query = `
		INSERT INTO api_keys (api_key, user_id, permissions, created_at)
		VALUES ($1::uuid, $2, 'READ_ONLY', NOW())
		RETURNING api_key::text
	`

	var insertedKey string
	err = r.Client.QueryRowContext(ctx, query, apiKeyUUID, userID).Scan(&insertedKey)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to insert new API key",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return "", wrapper.Wrap(err)
	}

	return insertedKey, nil
}

// GetAPIKeyByUserID возвращает последний активный ключ пользователя
func (r *Repository) GetAPIKeyByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	const query = `
		SELECT api_key::text 
		FROM api_keys 
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC 
		LIMIT 1
	`

	var apiKey string
	err := r.Client.GetContext(ctx, &apiKey, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", wrapper.Wrap(errors.New("no active API key found"))
		}
		r.Logger.ErrorContext(ctx, "failed to get API key by user",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return "", wrapper.Wrap(err)
	}

	return apiKey, nil
}

// GetUserIDByAPIKey находит user_id по API-ключу
func (r *Repository) GetUserIDByAPIKey(ctx context.Context, apiKey string) (uuid.UUID, error) {
	const query = `
		SELECT user_id 
		FROM api_keys 
		WHERE api_key::text = $1 AND is_active = true
	`

	var userID uuid.UUID
	err := r.Client.GetContext(ctx, &userID, query, apiKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, wrapper.Wrap(errors.New("api key not found or inactive"))
		}
		r.Logger.ErrorContext(ctx, "failed to get user by API key",
			slog.String("api_key", apiKey),
			slog.Any("error", err))
		return uuid.Nil, wrapper.Wrap(err)
	}

	return userID, nil
}
