// Файл: models/click.go
package models

import "time"

// Click представляет данные о клике пользователя, как они будут храниться в БД
type Click struct {
	ID 		  int 	  `json:"id"`
	UserID 	  int 	  `json:"user_id"` // Кто сделал клик
	OfferID   int 	  `json:"offer_id"`
	SubID 	  string  `json:"sub_id"`  // Уникальный идентификатор клика
	IPAddress string  `json:"ip_address"`
	UserAgent string  `json:"user_agent"`
	CreatedAt time.Time `json:"created_at"`
}

// ClickRequest DTO для приема минимальных данных от клиента
type ClickRequest struct {
	OfferID int    `json:"offer_id"`
	SubID   string `json:"sub_id"` 
}