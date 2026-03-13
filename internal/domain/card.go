package domain

import (
	"time"
)

type CardType string

const (
	CardTypeSubscriptions CardType = "subscriptions"
	CardTypeTravel        CardType = "travel"
	CardTypePremium       CardType = "premium"
)

type Card struct {
	ID               UUID       `json:"id"`
	UserID           UUID       `json:"userId"`
	ProviderCardID   string     `json:"providerCardId"`
	Bin              string     `json:"bin"`
	Last4Digits      string     `json:"last4Digits"`
	CardStatus       string     `json:"cardStatus"`
	Nickname         string     `json:"nickname,omitempty"`
	DailySpendLimit  Numeric    `json:"dailySpendLimit"`
	FailedAuthCount  int64      `json:"failedAuthCount"`
	CardType         CardType   `json:"cardType"`
	AutoTopUpEnabled bool       `json:"autoTopupTnabled"`
	AutoTopUpBelow   Numeric    `json:"autoTopupBelow"`
	AutoTopUpAmount  Numeric    `json:"autoTopupAmount"`
	Balance          Numeric    `json:"balance"`
	ExpiryDate       *time.Time `json:"expiryDate,omitempty"`
	CreatedAt        time.Time  `json:"createdAt"`
}

func NewCard(userID UUID, cardType CardType, providerCardID string) (*Card, error) {
	if !isValidCardType(cardType) {
		return nil, NewInvalidInput("invalid card_type: must be subscriptions, travel or premium")
	}

	return &Card{
		ID:              NewUUID(),
		UserID:          userID,
		ProviderCardID:  providerCardID,
		Bin:             "424242",
		Last4Digits:     "0000",
		CardStatus:      "ACTIVE",
		CardType:        cardType,
		DailySpendLimit: NewNumeric(1000),
		AutoTopUpBelow:  NewNumeric(100),
		AutoTopUpAmount: NewNumeric(500),
		Balance:         NewNumeric(0),
		CreatedAt:       time.Now().UTC(),
	}, nil
}

func isValidCardType(t CardType) bool {
	return t == CardTypeSubscriptions || t == CardTypeTravel || t == CardTypePremium
}
