package referral

import (
	"time"

	"github.com/google/uuid"
)

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
