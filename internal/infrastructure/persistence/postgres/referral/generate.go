package referral

import (
	"crypto/rand"
	"strings"

	"github.com/google/uuid"
)

// generateReferralCode создаёт уникальный реферальный код для пользователя
func generateReferralCode(userID uuid.UUID) string {
	// Берём первые 8 символов UUID (без дефисов) + 8 случайных символов
	return "USER" + strings.ReplaceAll(userID.String(), "-", "")[:8] + "-" + generateRandomString(8)
}

// generateRandomString возвращает крипто-случайную строку
func generateRandomString(n int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random string: " + err.Error()) // для MVP ок
	}
	for i := range b {
		b[i] = charset[b[i]%byte(len(charset))]
	}
	return string(b)
}
