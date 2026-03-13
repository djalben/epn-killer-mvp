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
	UserID           UUID       `json:"user_id"`
	ProviderCardID   string     `json:"provider_card_id"`
	Bin              string     `json:"bin"`
	Last4Digits      string     `json:"last_4_digits"`
	CardStatus       string     `json:"card_status"`
	Nickname         string     `json:"nickname,omitempty"`
	DailySpendLimit  Numeric    `json:"daily_spend_limit"`
	FailedAuthCount  int64      `json:"failed_auth_count"`
	CardType         CardType   `json:"card_type"`
	AutoTopUpEnabled bool       `json:"auto_topup_enabled"`
	AutoTopUpBelow   Numeric    `json:"auto_topup_below"`
	AutoTopUpAmount  Numeric    `json:"auto_topup_amount"`
	Balance          Numeric    `json:"balance"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
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
