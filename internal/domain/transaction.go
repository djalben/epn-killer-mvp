package domain

import (
	"time"
)

type Transaction struct {
	ID              UUID      `json:"id"`
	UserID          UUID      `json:"userId"`
	CardID          *UUID     `json:"cardId,omitempty"`
	Amount          Numeric   `json:"amount"`
	Fee             Numeric   `json:"fee"`
	TransactionType string    `json:"transactionType"`
	Status          string    `json:"status"`
	Details         string    `json:"details"`
	ProviderTxID    string    `json:"providerTxId,omitempty"`
	ExecutedAt      time.Time `json:"executedAt"`
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
