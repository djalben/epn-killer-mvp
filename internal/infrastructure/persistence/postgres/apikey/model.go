package apikey

import (
	"time"

	"github.com/google/uuid"
)

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
