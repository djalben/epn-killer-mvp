package domain

import (
	"time"
)

type Wallet struct {
	ID        UUID      `json:"id"`
	UserID    UUID      `json:"user_id"`
	Balance   Numeric   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewWallet(userID UUID) *Wallet {
	return &Wallet{
		ID:        NewUUID(),
		UserID:    userID,
		Balance:   NewNumeric(0),
		CreatedAt: time.Now().UTC(),
	}
}

func (w *Wallet) TopUp(amount Numeric) error {
	if amount.LessThanOrEqual(NewNumeric(0)) {
		return NewInvalidInput("amount must be positive")
	}
	w.Balance = w.Balance.Add(amount)
	return nil
}

func (w *Wallet) Withdraw(amount Numeric) error {
	if amount.LessThanOrEqual(NewNumeric(0)) {
		return NewInvalidInput("amount must be positive")
	}
	if w.Balance.LessThan(amount) {
		return NewInsufficientFunds()
	}
	w.Balance = w.Balance.Sub(amount)
	return nil
}
