package transaction

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gitlab.com/libs-artifex/wrapper/v2"

	"github.com/djalben/epn-killer-mvp/internal/entity"
)

// GetUserTransactionReport — извлекает отчёт по транзакциям пользователя с фильтрами
func (r *Repository) GetUserTransactionReport(ctx context.Context, userID uuid.UUID, filters map[string]interface{}) (*ReportSummary, error) {
	baseQuery := `
		SELECT id, user_id, card_id, amount, fee, transaction_type, status, details, executed_at
		FROM transactions
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIndex := 2

	// Применяем фильтры
	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		baseQuery += fmt.Sprintf(" AND executed_at >= $%d", argIndex)
		args = append(args, startDate)
		argIndex++
	}
	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		baseQuery += fmt.Sprintf(" AND executed_at <= $%d", argIndex)
		args = append(args, endDate)
		argIndex++
	}
	if txType, ok := filters["transaction_type"].(string); ok && txType != "" {
		baseQuery += fmt.Sprintf(" AND transaction_type = $%d", argIndex)
		args = append(args, txType)
		argIndex++
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}
	if cardID, ok := filters["card_id"].(uuid.UUID); ok && cardID != uuid.Nil {
		baseQuery += fmt.Sprintf(" AND card_id = $%d", argIndex)
		args = append(args, cardID)
		argIndex++
	}

	baseQuery += " ORDER BY executed_at DESC"

	rows, err := r.Client.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		r.Logger.ErrorContext(ctx, "failed to fetch user transactions",
			slog.String("user_id", userID.String()),
			slog.Any("filters", filters),
			slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}
	defer rows.Close()

	report := &ReportSummary{
		Transactions: make([]entity.Transaction, 0),
		TotalAmount:  decimal.Zero,
		TotalFee:     decimal.Zero,
	}

	for rows.Next() {
		var m transactionModel

		err := rows.Scan(
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
			r.Logger.ErrorContext(ctx, "failed to scan transaction row", slog.Any("error", err))
			continue
		}

		// Конвертируем внутреннюю модель в entity-модель для API
		report.Transactions = append(report.Transactions, m.toEntity())

		report.TotalTransactions++
		report.TotalAmount = report.TotalAmount.Add(decimal.NewFromFloat(m.Amount))
		report.TotalFee = report.TotalFee.Add(decimal.NewFromFloat(m.Fee))
	}

	if err := rows.Err(); err != nil {
		r.Logger.ErrorContext(ctx, "error iterating transaction rows", slog.Any("error", err))
		return nil, wrapper.Wrap(err)
	}

	r.Logger.InfoContext(ctx, "user transaction report fetched",
		slog.String("user_id", userID.String()),
		slog.Int("total_transactions", report.TotalTransactions))

	return report, nil
}
