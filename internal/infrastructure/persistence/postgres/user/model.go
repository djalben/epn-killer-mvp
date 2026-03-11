package user

import (
	"time"

	"github.com/google/uuid"
)

// userModel — пользователь системы
type userModel struct {
	ID             uuid.UUID `db:"id"`               // UUID из gen_random_uuid()
	Email          string    `db:"email"`            // Уникальный email
	PasswordHash   string    `db:"password_hash"`    // Хэш пароля
	Balance        float64   `db:"balance"`          // Основной баланс (Just-in-Time)
	Status         string    `db:"status"`           // ACTIVE, BLOCKED, PENDING_KYC и т.д.
	TelegramChatID *int64     `db:"telegram_chat_id"` // NULLable
	CreatedAt      time.Time `db:"created_at"`       // Автоматически NOW()
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
