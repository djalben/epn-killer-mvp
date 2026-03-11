package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

// --- СТРУКТУРЫ ЗАПРОСОВ К API (НОВЫЙ БЛОК) ---

// RegisterRequest - Запрос на регистрацию пользователя
type RegisterRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	ReferralCode string `json:"referral_code,omitempty"` // Опционально
}

// LoginRequest - Запрос на вход пользователя
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// FundRequest - Запрос на пополнение баланса (используется в deposit.go)
type FundRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

// AuthRequest - Запрос на авторизацию карты (используется в api/test_authorize.go)
type AuthRequest struct {
	CardID       string          `json:"card_id"`
	Amount       decimal.Decimal `json:"amount"`
	MerchantName string          `json:"merchant_name"`
}

// --- СТРУКТУРЫ ПОЛЬЗОВАТЕЛЕЙ И АУТЕНТИФИКАЦИИ ---

// User - Структура для пользователя
type User struct {
	ID             string          `json:"id"`
	Email          string          `json:"email"`
	PasswordHash   string          `json:"-"`
	Balance        decimal.Decimal `json:"balance"`
	CreatedAt      time.Time       `json:"created_at"`
	Status         string          `json:"status"`
	TeamID         *string         `json:"team_id,omitempty"` // или *uuid.UUID, если в entity UUID — строка
	TelegramChatID *int64          `json:"telegram_chat_id,omitempty"`
}

// APIKey - Структура для ключей
type APIKey struct {
	Key         string    `json:"api_key"`
	UserID      string    `json:"user_id"`
	Permissions string    `json:"permissions"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

// DepositRequest - Структура для запроса пополнения (используется API Key)
type DepositRequest struct {
	Amount decimal.Decimal `json:"amount"`
	UserID string          `json:"user_id"`
}

// --- СТРУКТУРЫ ТРАНЗАКЦИЙ И ОТЧЕТНОСТИ ---

// Transaction - Структура для транзакции
type Transaction struct {
	TransactionID   string          `json:"transaction_id"`
	UserID          string          `json:"user_id"`
	UserEmail       string          `json:"user_email,omitempty"`
	CardID          *string         `json:"card_id,omitempty"`
	CardLast4Digits string          `json:"card_last_4_digits,omitempty"`
	Amount          decimal.Decimal `json:"amount"`
	Fee             decimal.Decimal `json:"fee"`
	TransactionType string          `json:"transaction_type"`
	Status          string          `json:"status"`
	Details         string          `json:"details"`
	ExecutedAt      time.Time       `json:"executed_at"`
}

// --- СТРУКТУРЫ АВТОРИЗАЦИИ ---

// AuthResponse - Ответ от логики authorizeCard
type AuthResponse struct {
	Success bool            `json:"success"`
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Fee     decimal.Decimal `json:"fee"`
}

// --- СТРУКТУРЫ УПРАВЛЕНИЯ КАРТАМИ ---

// Card - Структура для карты
type Card struct {
	ID                     string          `json:"id"`
	UserID                 string          `json:"user_id"`
	TeamID                 *string         `json:"team_id,omitempty"` // Опционально, может быть NULL
	ProviderCardID         string          `json:"provider_card_id"`
	BIN                    string          `json:"bin"`
	Last4Digits            string          `json:"last_4_digits"`
	CardStatus             string          `json:"card_status"`
	Nickname               string          `json:"nickname"`
	DailySpendLimit        decimal.Decimal `json:"daily_spend_limit"`
	FailedAuthCount        int             `json:"failed_auth_count"`
	CardType               string          `json:"card_type"`
	AutoReplenishEnabled   bool            `json:"auto_replenish_enabled"`
	AutoReplenishThreshold decimal.Decimal `json:"auto_replenish_threshold"`
	AutoReplenishAmount    decimal.Decimal `json:"auto_replenish_amount"`
	CardBalance            decimal.Decimal `json:"card_balance"` // Текущий баланс карты
	CreatedAt              time.Time       `json:"created_at"`
}

// MassIssueRequest - Запрос на массовый выпуск карт
type MassIssueRequest struct {
	Count        int             `json:"count"`
	DailyLimit   decimal.Decimal `json:"daily_limit"`
	CardNickname string          `json:"nickname"`
	MerchantName string          `json:"merchant_name"`
	CardType     string          `json:"card_type"`         // VISA или MasterCard
	TeamID       *string         `json:"team_id,omitempty"` // Опционально, для команд
}

// CardIssueResult - Результат выпуска одной карты
type CardIssueResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Card      *Card  `json:"card,omitempty"`
	Status    string `json:"status"`
	CardLast4 string `json:"card_last_4"`
	Nickname  string `json:"nickname"`
}

// MassIssueResponse - Ответ на массовый выпуск карт
type MassIssueResponse struct {
	Successful int               `json:"successful_count"`
	Failed     int               `json:"failed_count"`
	Results    []CardIssueResult `json:"results"`
}

// --- СТРУКТУРЫ АВТОПОПОЛНЕНИЯ КАРТ ---

// AutoReplenishRequest - Запрос на настройку автопополнения
type AutoReplenishRequest struct {
	Enabled   bool            `json:"enabled"`
	Threshold decimal.Decimal `json:"threshold"`
	Amount    decimal.Decimal `json:"amount"`
}

// --- СТРУКТУРЫ КОМАНД ---

// Team - Команда
type Team struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TeamMember - Участник команды
type TeamMember struct {
	ID        string    `json:"id"`
	TeamID    string    `json:"team_id"`
	UserID    string    `json:"user_id"`
	Role      string    `json:"role"` // 'owner', 'admin', 'member'
	InvitedBy *string   `json:"invited_by,omitempty"`
	JoinedAt  time.Time `json:"joined_at"`
	User      *User     `json:"user,omitempty"` // Для деталей пользователя
}

// CreateTeamRequest - Запрос на создание команды
type CreateTeamRequest struct {
	Name string `json:"name"`
}

// InviteTeamMemberRequest - Запрос на приглашение участника
type InviteTeamMemberRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"` // 'admin' или 'member'
}

// UpdateTeamMemberRoleRequest - Запрос на изменение роли
type UpdateTeamMemberRoleRequest struct {
	Role string `json:"role"`
}

// --- СТРУКТУРЫ GRADE СИСТЕМЫ ---

// UserGrade - Grade пользователя
type UserGrade struct {
	ID         string          `json:"id"`
	UserID     string          `json:"user_id"`
	Grade      string          `json:"grade"` // 'STANDARD', 'SILVER', 'GOLD', 'PLATINUM', 'BLACK'
	TotalSpent decimal.Decimal `json:"total_spent"`
	FeePercent decimal.Decimal `json:"fee_percent"` // Комиссия в процентах (6.70 = 6.7%)
	UpdatedAt  time.Time       `json:"updated_at"`
}

// GradeInfo - Информация о Grade для отображения
type GradeInfo struct {
	Grade      string           `json:"grade"`
	TotalSpent decimal.Decimal  `json:"total_spent"`
	FeePercent decimal.Decimal  `json:"fee_percent"`
	NextGrade  *string          `json:"next_grade,omitempty"`
	NextSpend  *decimal.Decimal `json:"next_spend,omitempty"` // Сколько нужно потратить до следующего уровня
}

// --- СТРУКТУРЫ РЕФЕРАЛЬНОЙ ПРОГРАММЫ ---

// Referral - Реферал
type Referral struct {
	ID               string          `json:"id"`
	ReferrerID       string          `json:"referrer_id"`
	ReferredID       string          `json:"referred_id"`
	ReferralCode     string          `json:"referral_code"`
	Status           string          `json:"status"` // 'PENDING', 'ACTIVE', 'COMPLETED'
	CommissionEarned decimal.Decimal `json:"commission_earned"`
	CreatedAt        time.Time       `json:"created_at"`
}

// ReferralStats - Статистика реферальной программы
type ReferralStats struct {
	TotalReferrals  int             `json:"total_referrals"`
	ActiveReferrals int             `json:"active_referrals"`
	TotalCommission decimal.Decimal `json:"total_commission"`
	ReferralCode    string          `json:"referral_code"`
}

// ReportSummary - Агрегированная сводка по кликам
// Это главная структура, которую мы будем возвращать пользователю (UI-отчет).
type ReportSummary struct {
	Date         time.Time `json:"date"`          // Дата, по которой сгруппированы данные
	SubID        string    `json:"sub_id"`        // SubID, по которому сгруппированы данные
	TotalClicks  int       `json:"total_clicks"`  // Общее количество кликов за период
	UniqueClicks int       `json:"unique_clicks"` // Количество уникальных кликов (по IP)
	// Spend - Будет добавлено, когда реализуем финансовую логику (ЭТАП 2)
	// Conversions - Будет добавлено, когда реализуем логику конверсий
}

// SpendReportSummary - Агрегированная сводка по расходам (для трекеров)
// Используется для API-интеграции (GET /api/v1/data/spends)
type SpendReportSummary struct {
	Date             time.Time `json:"date"`              // Дата, по которой сгруппированы данные
	SubID            string    `json:"sub_id"`            // SubID, по которому сгруппированы данные
	TotalSpend       float64   `json:"total_spend"`       // Общая сумма трат за период (в USD)
	TransactionCount int       `json:"transaction_count"` // Общее количество транзакций
}

// ReportRequest - Структура для парсинга параметров запроса отчета
type ReportRequest struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	SubID     string    `json:"sub_id"`   // Опциональный фильтр
	GroupBy   string    `json:"group_by"` // 'day', 'subid', 'offer'
}
