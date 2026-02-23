package transaction

import (
	"time"

	"github.com/djalben/epn-killer-mvp/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// transaction — любая транзакция
type transactionModel struct {
	ID              uuid.UUID  `db:"id"`
	UserID          uuid.UUID  `db:"user_id"`
	CardID          *uuid.UUID `db:"card_id"` // Может быть NULL для пополнений
	Amount          float64    `db:"amount"`
	Fee             float64    `db:"fee"`
	TransactionType string     `db:"transaction_type"` // FUND, AUTH, CAPTURE, DECLINE и т.д.
	Status          string     `db:"status"`
	Details         *string    `db:"details"` // Merchant name, описание
	ExecutedAt      time.Time  `db:"executed_at"`
}

// ReportSummary - Агрегированная сводка по транзакциям пользователя (для API)
type ReportSummary struct {
	TotalTransactions int                  `json:"total_transactions"`
	TotalAmount       decimal.Decimal      `json:"total_amount"`
	TotalFee          decimal.Decimal      `json:"total_fee"`
	Transactions      []entity.Transaction `json:"transactions"`
}

// toEntity конвертирует внутреннюю модель транзакции в entity-модель для API
func (m transactionModel) toEntity() entity.Transaction {
	tx := entity.Transaction{
		TransactionID:   m.ID.String(),
		UserID:          m.UserID.String(),
		Amount:          decimal.NewFromFloat(m.Amount),
		Fee:             decimal.NewFromFloat(m.Fee),
		TransactionType: m.TransactionType,
		Status:          m.Status,
		Details:         "",
		ExecutedAt:      m.ExecutedAt,
	}

	if m.CardID != nil {
		cardIDStr := m.CardID.String()
		tx.CardID = &cardIDStr
	}

	if m.Details != nil {
		tx.Details = *m.Details
	}

	return tx
}
