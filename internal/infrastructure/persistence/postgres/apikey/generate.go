package apikey

import (
	"fmt"

	"crypto/rand"
	"encoding/hex"

	"gitlab.com/libs-artifex/wrapper/v2"
)

// generateRandomString генерирует случайную строку в шестнадцатеричном формате
func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", wrapper.Wrap(err)
	}
	return hex.EncodeToString(b), nil
}

// formatAsUUID форматирует 32 hex-символа в вид UUID
func formatAsUUID(hexStr string) string {
	if len(hexStr) != 32 {
		return hexStr
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hexStr[0:8], hexStr[8:12], hexStr[12:16], hexStr[16:20], hexStr[20:32])
}
