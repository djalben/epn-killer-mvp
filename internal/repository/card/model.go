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

// apiKey — ключ для интеграции с трекерами
type apiKeyModel struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	APIKey      uuid.UUID `db:"api_key"`     // UUID v4
	Permissions string    `db:"permissions"` // READ_ONLY
	Description *string   `db:"description"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
}

// team — команда
type teamModel struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	OwnerID   uuid.UUID `db:"owner_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// teamMember — участник команды
type teamMemberModel struct {
	ID        uuid.UUID  `db:"id"`
	TeamID    uuid.UUID  `db:"team_id"`
	UserID    uuid.UUID  `db:"user_id"`
	Role      string     `db:"role"` // owner, admin, member
	InvitedBy *uuid.UUID `db:"invited_by"`
	JoinedAt  time.Time  `db:"joined_at"`
}

// userGrade — уровень пользователя
type userGradeModel struct {
	ID         uuid.UUID `db:"id"`
	UserID     uuid.UUID `db:"user_id"` // UNIQUE
	Grade      string    `db:"grade"`
	TotalSpent float64   `db:"total_spent"`
	FeePercent float64   `db:"fee_percent"` // 6.70 для STANDARD
	UpdatedAt  time.Time `db:"updated_at"`
}

// referral — реферальная связь
type referralModel struct {
	ID               uuid.UUID `db:"id"`
	ReferrerID       uuid.UUID `db:"referrer_id"`
	ReferredID       uuid.UUID `db:"referred_id"`
	ReferralCode     string    `db:"referral_code"` // UNIQUE
	Status           string    `db:"status"`        // PENDING, ACTIVE, COMPLETED
	CommissionEarned float64   `db:"commission_earned"`
	CreatedAt        time.Time `db:"created_at"`
}
