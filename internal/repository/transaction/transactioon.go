package transaction

import (
	"fmt"
	"log"
	"time"

	"github.com/djalben/epn-killer-mvp/internal/notification"
	"github.com/shopspring/decimal"
)

// ProcessDeposit - Обрабатывает пополнение баланса пользователя и записывает транзакцию.
func ProcessDeposit(userID int, amount decimal.Decimal) error {
	if GlobalDB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("сумма пополнения должна быть положительной")
	}

	// 1. НАЧАЛО ТРАНЗАКЦИИ
	tx, err := GlobalDB.Begin()
	if err != nil {
		log.Printf("DB Error Begin: %v", err)
		return fmt.Errorf("не удалось начать транзакцию")
	}
	defer tx.Rollback()

	// 2. Увеличение баланса пользователя (атомарно)
	_, err = tx.Exec(
		"UPDATE users SET balance = balance + $1 WHERE id = $2",
		amount, userID,
	)
	if err != nil {
		log.Printf("DB Error Update Balance: %v", err)
		return fmt.Errorf("не удалось обновить баланс")
	}

	// 3. Запись транзакции пополнения (FUND)
	_, err = tx.Exec(
		`INSERT INTO transactions (user_id, amount, fee, transaction_type, status, details, executed_at)
			VALUES ($1, $2, $3, 'FUND', 'APPROVED', $4, $5)`,
		userID,
		amount,
		decimal.Zero, // Комиссия
		fmt.Sprintf("Deposit via API. Amount: %s", amount.String()),
		time.Now(),
	)
	if err != nil {
		log.Printf("DB Error Insert Transaction: %v", err)
		return fmt.Errorf("не удалось записать транзакцию")
	}

	// 4. КОММИТ
	if err := tx.Commit(); err != nil {
		log.Printf("DB Error Commit: %v", err)
		return fmt.Errorf("ошибка фиксации транзакции")
	}

	log.Printf("User %d successfully deposited %s. Transaction committed.", userID, amount.String())

	// 5. ОТПРАВКА TELEGRAM УВЕДОМЛЕНИЯ
	user, err := GetUserByID(userID)
	if err == nil && user.TelegramChatID.Valid {
		message := fmt.Sprintf("✅ Deposit Successful!\n\nAmount: $%s\nNew Balance: $%s",
			amount.String(),
			user.Balance.String())
		notification.SendTelegramMessage(user.TelegramChatID.Int64, message)
	}

	return nil
}
