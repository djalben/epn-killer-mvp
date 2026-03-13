package domain

import (
	"time"
)

type User struct {
	ID             UUID       `json:"id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"`
	KYCStatus      KYCStatus  `json:"kyc_status"`
	Status         UserStatus `json:"status"`
	TelegramChatID *int64     `json:"telegram_chat_id,omitempty"`
	ReferralCode   string     `json:"referral_code"`
	ReferredBy     *UUID      `json:"referred_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type UserStatus string
type KYCStatus string

const (
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusBlocked UserStatus = "BLOCKED"

	KYCPending  KYCStatus = "PENDING"
	KYCApproved KYCStatus = "APPROVED"
	KYCRejected KYCStatus = "REJECTED"
)

func NewUser(email, passwordHash string) (*User, error) {
	if email == "" {
		return nil, NewInvalidInput("email is required")
	}
	return &User{
		ID:           NewUUID(),
		Email:        email,
		PasswordHash: passwordHash,
		KYCStatus:    KYCPending,
		Status:       UserStatusActive,
		CreatedAt:    time.Now().UTC(),
	}, nil
}
