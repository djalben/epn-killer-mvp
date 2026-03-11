package referral

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.com/libs-artifex/wrapper/v2"

	"github.com/djalben/epn-killer-mvp/internal/entity"
)

// ──────────────────────────────────────────────────────────────
// Основные методы репозитория
// ──────────────────────────────────────────────────────────────

// CreateReferral создаёт реферальную запись
func (r *Repository) CreateReferral(ctx context.Context, referrerID, referredID uuid.UUID, referralCode string) error {
	const query = `
		INSERT INTO referrals (referrer_id, referred_id, referral_code, status)
		VALUES ($1, $2, $3, 'ACTIVE')
	`

	_, err := r.Client.ExecContext(ctx, query, referrerID, referredID, referralCode)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to create referral",
			slog.String("referrer_id", referrerID.String()),
			slog.String("referred_id", referredID.String()),
			slog.String("referral_code", referralCode),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "referral created",
		slog.String("referrer_id", referrerID.String()),
		slog.String("referred_id", referredID.String()),
		slog.String("referral_code", referralCode))

	return nil
}

// GetUserReferralCode получает реферальный код пользователя (создаёт новый, если нет)
func (r *Repository) GetUserReferralCode(ctx context.Context, userID uuid.UUID) (string, error) {
	const selectQuery = `
		SELECT referral_code 
		FROM referrals 
		WHERE referrer_id = $1 
		LIMIT 1
	`

	var code string
	err := r.Client.GetContext(ctx, &code, selectQuery, userID)
	if err == nil {
		return code, nil
	}

	if !errors.Is(err, sql.ErrNoRows) {
		r.Logger.ErrorContext(ctx, "failed to get referral code",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return "", wrapper.Wrap(err)
	}

	// Кода нет — генерируем новый
	newCode := generateReferralCode(userID)

	// Проверяем уникальность
	for {
		var existing string
		err = r.Client.GetContext(ctx, &existing, "SELECT referral_code FROM referrals WHERE referral_code = $1", newCode)
		if errors.Is(err, sql.ErrNoRows) {
			break // код уникальный
		}
		if err != nil {
			r.Logger.ErrorContext(ctx, "failed to check referral code uniqueness",
				slog.String("user_id", userID.String()),
				slog.Any("error", err))
			return "", wrapper.Wrap(err)
		}
		newCode = generateReferralCode(userID)
	}

	// Создаём запись
	err = r.CreateReferral(ctx, userID, uuid.Nil, newCode)
	if err != nil {
		return "", wrapper.Wrap(err)
	}

	return newCode, nil
}

// GetReferralStats получает статистику реферальной программы
func (r *Repository) GetReferralStats(ctx context.Context, userID uuid.UUID) (*entity.ReferralStats, error) {
	code, err := r.GetUserReferralCode(ctx, userID)
	if err != nil {
		return nil, wrapper.Wrap(err)
	}

	const statsQuery = `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'ACTIVE' THEN 1 END) as active,
			COALESCE(SUM(commission_earned), 0) as commission
		FROM referrals 
		WHERE referrer_id = $1
	`

	var stats struct {
		Total      int             `db:"total"`
		Active     int             `db:"active"`
		Commission decimal.Decimal `db:"commission"`
	}

	err = r.Client.GetContext(ctx, &stats, statsQuery, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to get referral stats",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	return &entity.ReferralStats{
		TotalReferrals:  stats.Total,
		ActiveReferrals: stats.Active,
		TotalCommission: stats.Commission,
		ReferralCode:    code,
	}, nil
}

// ProcessReferralRegistration обрабатывает регистрацию по реферальной ссылке
func (r *Repository) ProcessReferralRegistration(ctx context.Context, referredID uuid.UUID, referralCode string) error {
	const findReferrerQuery = `
		SELECT referrer_id 
		FROM referrals 
		WHERE referral_code = $1 AND status = 'ACTIVE' 
		LIMIT 1
	`

	var referrerID uuid.UUID
	err := r.Client.GetContext(ctx, &referrerID, findReferrerQuery, referralCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wrapper.Wrap(errors.New("invalid or inactive referral code"))
		}
		r.Logger.ErrorContext(ctx, "failed to find referrer by code",
			slog.String("referral_code", referralCode),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	// Проверяем, не зарегистрирован ли уже
	const checkExistsQuery = `
		SELECT id 
		FROM referrals 
		WHERE referrer_id = $1 AND referred_id = $2
		LIMIT 1
	`

	var existingID uuid.UUID
	err = r.Client.GetContext(ctx, &existingID, checkExistsQuery, referrerID, referredID)
	if err == nil {
		return nil // уже существует — не ошибка
	}
	if !errors.Is(err, sql.ErrNoRows) {
		r.Logger.ErrorContext(ctx, "failed to check existing referral",
			slog.String("referrer_id", referrerID.String()),
			slog.String("referred_id", referredID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	return r.CreateReferral(ctx, referrerID, referredID, referralCode)
}

// AddReferralCommission добавляет комиссию рефереру
func (r *Repository) AddReferralCommission(ctx context.Context, referrerID uuid.UUID, amount decimal.Decimal) error {
	const updateQuery = `
		UPDATE referrals 
		SET commission_earned = commission_earned + $1 
		WHERE referrer_id = $2 AND status = 'ACTIVE'
	`

	_, err := r.Client.ExecContext(ctx, updateQuery, amount, referrerID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to add referral commission",
			slog.String("referrer_id", referrerID.String()),
			slog.Any("amount", amount),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	return nil
}
