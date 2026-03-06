-- +goose Up
-- +goose StatementBegin
-- UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- 1. Таблица пользователей
CREATE TABLE users (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    -- Баланс хранится только здесь (Just-in-Time Funding)
    balance NUMERIC(20, 4) DEFAULT 0.0000 NOT NULL,
    status VARCHAR(50) DEFAULT 'ACTIVE',
    telegram_chat_id BIGINT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 2. Таблица команд
CREATE TABLE teams (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. Таблица карт
CREATE TABLE cards (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    provider_card_id VARCHAR(100) NOT NULL,
    bin VARCHAR(6) NOT NULL DEFAULT '424242',
    last_4_digits VARCHAR(4) NOT NULL,
    card_status VARCHAR(50) DEFAULT 'ACTIVE',
    nickname VARCHAR(100),
    service_slug VARCHAR(50) DEFAULT 'arbitrage',
    daily_spend_limit NUMERIC(20, 4) DEFAULT 1000.0000,
    failed_auth_count BIGINT DEFAULT 0,
    card_type VARCHAR(20) DEFAULT 'VISA',
    auto_replenish_enabled BOOLEAN DEFAULT FALSE,
    auto_replenish_threshold NUMERIC(20, 4) DEFAULT 0.0000,
    auto_replenish_amount NUMERIC(20, 4) DEFAULT 0.0000,
    card_balance NUMERIC(20, 4) DEFAULT 0.0000,
    team_id UUID REFERENCES teams(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 4. Таблица транзакций
CREATE TABLE transactions (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    card_id UUID REFERENCES cards(id),
    amount NUMERIC(20, 4) NOT NULL,
    fee NUMERIC(20, 4) DEFAULT 0.0000,
    transaction_type VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    details TEXT,
    provider_tx_id VARCHAR(255),
    executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 5. API Ключи (Для трекеров)
CREATE TABLE api_keys (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    api_key UUID UNIQUE DEFAULT gen_random_uuid(),
    permissions VARCHAR(50) DEFAULT 'READ_ONLY', -- READ_ONLY для Keitaro/Binom
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 6. Таблица участников команд
CREATE TABLE team_members (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    team_id UUID REFERENCES teams(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member', -- 'owner', 'admin', 'member'
    invited_by UUID REFERENCES users(id),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT team_id_user_id_unq UNIQUE(team_id, user_id)
);

-- 7. Таблица Grade пользователей (система уровней)
CREATE TABLE user_grades (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    grade VARCHAR(50) DEFAULT 'STANDARD', -- 'STANDARD', 'SILVER', 'GOLD', 'PLATINUM', 'BLACK'
    total_spent NUMERIC(20, 4) DEFAULT 0.0000, -- Общая сумма трат
    fee_percent NUMERIC(5, 2) DEFAULT 6.70, -- Комиссия в процентах (6.7% для STANDARD)
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 8. Таблица реферальной программы
CREATE TABLE referrals (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    referrer_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Кто пригласил
    referred_id UUID REFERENCES users(id) ON DELETE CASCADE, -- Кого пригласили
    referral_code VARCHAR(50) UNIQUE NOT NULL, -- Уникальный код реферала
    status VARCHAR(50) DEFAULT 'PENDING', -- 'PENDING', 'ACTIVE', 'COMPLETED'
    commission_earned NUMERIC(20, 4) DEFAULT 0.0000, -- Заработанная комиссия
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 9. Курсы валют (для конвертаций и Wallester)
CREATE TABLE exchange_rates (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    currency_from VARCHAR(10) NOT NULL,
    currency_to VARCHAR(10) NOT NULL,
    base_rate NUMERIC(20, 8) NOT NULL,
    markup_percent NUMERIC(6, 2) NOT NULL DEFAULT 0.00,
    final_rate NUMERIC(20, 8) NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT exchange_rates_currency_pair_unq UNIQUE(currency_from, currency_to)
);

-- 10. Services (каталоги для Wallester)
CREATE TABLE services (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
INSERT INTO services (slug, name) VALUES
    ('arbitrage', 'Arbitrage'),
    ('travel', 'Travel'),
    ('services', 'Services'),
    ('subscriptions', 'Subscriptions');

-- Индексы для скорости (после всех таблиц)
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_cards_user_id ON cards(user_id);
CREATE INDEX idx_cards_team_id ON cards(team_id);
CREATE INDEX idx_cards_user_status ON cards(user_id, card_status);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_card_id ON transactions(card_id);
CREATE INDEX idx_transactions_executed_at ON transactions(executed_at DESC);
CREATE INDEX idx_transactions_provider_tx_id ON transactions(provider_tx_id) WHERE provider_tx_id IS NOT NULL;
CREATE INDEX idx_teams_owner_id ON teams(owner_id);
CREATE INDEX idx_team_members_team_id ON team_members(team_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);
CREATE INDEX idx_user_grades_user_id ON user_grades(user_id);
CREATE INDEX idx_referrals_referrer_id ON referrals(referrer_id);
CREATE INDEX idx_referrals_referred_id ON referrals(referred_id);
CREATE INDEX idx_referrals_code ON referrals(referral_code);
CREATE INDEX idx_exchange_rates_pair ON exchange_rates(currency_from, currency_to);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS exchange_rates;
DROP TABLE IF EXISTS referrals;
DROP TABLE IF EXISTS user_grades;
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
