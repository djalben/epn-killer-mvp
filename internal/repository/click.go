package repository

import (
	"fmt"
	"log"

	"github.com/djalben/epn-killer-mvp/internal/models"
)

// GlobalDB должна быть определена в main.go
// var GlobalDB *sql.DB

// SaveClick сохраняет данные о клике в базе данных и возвращает ID созданной записи.
func SaveClick(click models.Click) (int, error) { // ВАЖНО: Возвращает ID
	if GlobalDB == nil {
		return 0, fmt.Errorf("database connection not initialized")
	}

	// --- !!! ОКОНЧАТЕЛЬНО ИСПРАВЛЕНО: СТОЛБЕЦ 'status' ИСКЛЮЧЕН, 
    // чтобы заполнялся значением по умолчанию в БД и не нарушал NOT NULL !!! ---
	query := `
		INSERT INTO clicks (user_id, offer_id, sub_id, ip_address, user_agent, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	var clickID int
	
	err := GlobalDB.QueryRow(
		query,
		click.UserID,
		click.OfferID,
		click.SubID,
		click.IPAddress,
		click.UserAgent,
		click.CreatedAt, // <-- Шестая переменная ($6)
	).Scan(&clickID)

	if err != nil {
		log.Printf("Error saving click: %v", err)
		return 0, err
	}

	return clickID, nil
}