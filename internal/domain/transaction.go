package domain

import (
	"time"
)

type Transaction struct {
	ID              UUID      `json:"id"`
	UserID          UUID      `json:"user_id"`
	CardID          *UUID     `json:"card_id,omitempty"`
	Amount          Numeric   `json:"amount"`
	Fee             Numeric   `json:"fee"`
	TransactionType string    `json:"transaction_type"`
	Status          string    `json:"status"`
	Details         string    `json:"details"`
	ProviderTxID    string    `json:"provider_tx_id,omitempty"`
	ExecutedAt      time.Time `json:"executed_at"`
}

func NewTransaction(
	userID UUID,
	cardID *UUID,
	amount, fee Numeric,
	txType, status, details string,
) *Transaction {
	return &Transaction{
		ID:              NewUUID(),
		UserID:          userID,
		CardID:          cardID,
		Amount:          amount,
		Fee:             fee,
		TransactionType: txType,
		Status:          status,
		Details:         details,
		ExecutedAt:      time.Now().UTC(),
	}
}
