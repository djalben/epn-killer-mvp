package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// GlobalDB должен быть объявлен в этом пакете (например, globals.go)

// --- ФУНКЦИИ АУТЕНТИФИКАЦИИ И ПОЛУЧЕНИЯ ДАННЫХ ---

// CreateUser создает нового пользователя и сохраняет его в БД.
func CreateUser(user userModel) (userModel, error) {
	if GlobalDB == nil {
		return userModel{}, fmt.Errorf("database connection not initialized")
	}

	queryUser := `
		INSERT INTO users (email, password_hash, balance, status) 
		VALUES ($1, $2, 0.00, 'ACTIVE') 
		RETURNING id, created_at, balance
	`
	var createdUser userModel

	err := GlobalDB.QueryRow(queryUser, user.Email, user.PasswordHash).
		Scan(&createdUser.ID, &createdUser.CreatedAt, &createdUser.Balance)

	if err != nil {
		log.Printf("Error creating user %s: %v", user.Email, err)
		return userModel{}, err
	}

	createdUser.Email = user.Email

	log.Printf("User %s created with ID: %d", createdUser.Email, createdUser.ID)

	// Генерируем API-ключ для нового пользователя.
	// TODO: вынести это отдельно.
	_, err = GenerateAPIKey(createdUser.ID)
	if err != nil {
		log.Printf("WARNING: Could not generate API key for user %d on creation: %v", createdUser.ID, err)
	}

	return createdUser, nil
}

// GetUserByEmail - Находит пользователя по email.
func (r *Repository) GetUserByEmail(ctx context.Context, email string) (userModel, error) {
	if GlobalDB == nil {
		return userModel{}, fmt.Errorf("database connection not initialized")
	}

	// ИСПРАВЛЕНО: TelegramChatID сканируется напрямую в модель (sql.NullInt64)
	query := `SELECT id, email, password_hash, balance, created_at, telegram_chat_id FROM users WHERE email = $1`

	var user userModel

	err := GlobalDB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
		&user.TelegramChatID, // Сканируем sql.NullInt64
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userModel{}, errors.New("пользователь не найден")
		}
		log.Printf("DB Error GetUserByEmail: %v", err)
		return userModel{}, err
	}

	return user, nil
}

// GetUserByID - Находит пользователя по ID.
func GetUserByID(userID int) (userModel, error) {
	if GlobalDB == nil {
		return userModel{}, fmt.Errorf("database connection not initialized")
	}

	// ИСПРАВЛЕНО: TelegramChatID сканируется напрямую в модель (sql.NullInt64)
	query := `SELECT id, email, password_hash, balance, created_at, telegram_chat_id FROM users WHERE id = $1`

	var user userModel

	err := GlobalDB.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Balance,
		&user.CreatedAt,
		&user.TelegramChatID, // Сканируем sql.NullInt64
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userModel{}, errors.New("пользователь не найден")
		}
		log.Printf("DB Error GetUserByID: %v", err)
		return userModel{}, err
	}

	return user, nil
}

// UpdateTelegramChatID - Обновляет TelegramChatID пользователя.
// Принимает int, так как ChatID в Go - это int, а в БД - BIGINT.
func UpdateTelegramChatID(userID int, chatID int) error {
	if GlobalDB == nil {
		return fmt.Errorf("database connection not initialized")
	}

	_, err := GlobalDB.Exec(
		"UPDATE users SET telegram_chat_id = $1 WHERE id = $2",
		chatID, userID,
	)
	if err != nil {
		log.Printf("DB Error UpdateTelegramChatID: %v", err)
		return fmt.Errorf("не удалось обновить telegram_chat_id")
	}
	return nil
}
