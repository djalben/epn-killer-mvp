package card

import (
	"time"

	"github.com/google/uuid"
)

// cardModel — виртуальная карта
type cardModel struct {
	ID                     uuid.UUID  `db:"id"`
	UserID                 uuid.UUID  `db:"user_id"`
	ProviderCardID         string     `db:"provider_card_id"` // ID от Wallester
	BIN                    string     `db:"bin"`              // Первые 6 цифр
	Last4Digits            string     `db:"last_4_digits"`
	CardStatus             string     `db:"card_status"` // ACTIVE, BLOCKED и т.д.
	Nickname               *string    `db:"nickname"`    // Пользовательское имя карты
	DailySpendLimit        float64    `db:"daily_spend_limit"`
	FailedAuthCount        int64      `db:"failed_auth_count"`
	CardType               string     `db:"card_type"` // VISA, MASTERCARD
	AutoReplenishEnabled   bool       `db:"auto_replenish_enabled"`
	AutoReplenishThreshold float64    `db:"auto_replenish_threshold"`
	AutoReplenishAmount    float64    `db:"auto_replenish_amount"`
	CardBalance            float64    `db:"card_balance"` // Баланс карты
	TeamID                 *uuid.UUID `db:"team_id"`      // NULLable
	CreatedAt              time.Time  `db:"created_at"`
}
