package user

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

// Grade thresholds
const (
	GradeStandard = "STANDARD"
	GradeSilver   = "SILVER"
	GradeGold     = "GOLD"
	GradePlatinum = "PLATINUM"
	GradeBlack    = "BLACK"
)

// Grade fee percentages
var (
	FeeStandard = decimal.NewFromFloat(6.70)
	FeeSilver   = decimal.NewFromFloat(6.00)
	FeeGold     = decimal.NewFromFloat(5.00)
	FeePlatinum = decimal.NewFromFloat(4.00)
	FeeBlack    = decimal.NewFromFloat(3.00)
)

// GetUserGrade возвращает текущий Grade пользователя
func (r *Repository) GetUserGrade(ctx context.Context, userID uuid.UUID) (*entity.UserGrade, error) {
	const query = `
		SELECT id, user_id, grade, total_spent, fee_percent, updated_at 
		FROM user_grades 
		WHERE user_id = $1
	`

	var grade entity.UserGrade
	err := r.Client.GetContext(ctx, &grade, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return r.CreateUserGrade(ctx, userID)
		}
		r.Logger.ErrorContext(ctx, "failed to get user grade",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	return &grade, nil
}

// CreateUserGrade создаёт Grade для пользователя (по умолчанию STANDARD)
func (r *Repository) CreateUserGrade(ctx context.Context, userID uuid.UUID) (*entity.UserGrade, error) {
	const query = `
		INSERT INTO user_grades (user_id, grade, total_spent, fee_percent)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, grade, total_spent, fee_percent, updated_at
	`

	var grade entity.UserGrade
	err := r.Client.QueryRowContext(ctx, query, userID, GradeStandard, decimal.Zero, FeeStandard).
		Scan(&grade.ID, &grade.UserID, &grade.Grade, &grade.TotalSpent, &grade.FeePercent, &grade.UpdatedAt)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to create user grade",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	return &grade, nil
}

// CalculateGradeFromSpent вычисляет Grade и комиссию на основе общей суммы трат
func CalculateGradeFromSpent(totalSpent decimal.Decimal) (string, decimal.Decimal) {
	thresholdBlack := decimal.NewFromInt(100000)
	thresholdPlatinum := decimal.NewFromInt(50000)
	thresholdGold := decimal.NewFromInt(10000)
	thresholdSilver := decimal.NewFromInt(1000)

	switch {
	case totalSpent.GreaterThanOrEqual(thresholdBlack):
		return GradeBlack, FeeBlack
	case totalSpent.GreaterThanOrEqual(thresholdPlatinum):
		return GradePlatinum, FeePlatinum
	case totalSpent.GreaterThanOrEqual(thresholdGold):
		return GradeGold, FeeGold
	case totalSpent.GreaterThanOrEqual(thresholdSilver):
		return GradeSilver, FeeSilver
	default:
		return GradeStandard, FeeStandard
	}
}

// UpdateUserGrade обновляет Grade пользователя на основе его трат
func (r *Repository) UpdateUserGrade(ctx context.Context, userID uuid.UUID) error {
	// Считаем общую сумму APPROVED CAPTURE транзакций
	const totalQuery = `
		SELECT COALESCE(SUM(amount), 0) 
		FROM transactions 
		WHERE user_id = $1 
		  AND transaction_type = 'CAPTURE' 
		  AND status = 'APPROVED'
	`

	var totalSpent decimal.Decimal
	err := r.Client.GetContext(ctx, &totalSpent, totalQuery, userID)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to calculate total spent",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	newGrade, newFee := CalculateGradeFromSpent(totalSpent)

	const updateQuery = `
		INSERT INTO user_grades (user_id, grade, total_spent, fee_percent, updated_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			grade = $2,
			total_spent = $3,
			fee_percent = $4,
			updated_at = NOW()
	`

	_, err = r.Client.ExecContext(ctx, updateQuery, userID, newGrade, totalSpent, newFee)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to update user grade",
			slog.String("user_id", userID.String()),
			slog.Any("error", err))
		return wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "user grade updated",
		slog.String("user_id", userID.String()),
		slog.String("new_grade", newGrade),
		slog.Any("total_spent", totalSpent))

	return nil
}
