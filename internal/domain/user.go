package domain

import (
	"time"
)

type User struct {
	ID             UUID       `json:"id"`
	Email          string     `json:"email"`
	PasswordHash   string     `json:"-"`
	KYCStatus      KYCStatus  `json:"kycStatus"`
	Status         UserStatus `json:"status"`
	TelegramChatID *int64     `json:"telegramChatId,omitempty"`
	ReferralCode   string     `json:"referralCode"`
	ReferredBy     *UUID      `json:"referredBy,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
}

type (
	UserStatus string
	KYCStatus  string
)

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
