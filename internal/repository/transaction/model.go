package transaction

import (
	"time"

	"github.com/google/uuid"
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
