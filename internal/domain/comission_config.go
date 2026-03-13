package domain

import "time"

type CommissionConfig struct {
	ID          UUID      `json:"id"`
	Key         string    `json:"key"`
	Value       Numeric   `json:"value"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

const (
	FeeStandard     = "fee_standard"
	FeeSilver       = "fee_silver"
	FeeGold         = "fee_gold"
	FeePlatinum     = "fee_platinum"
	FeeBlack        = "fee_black"
	ReferralPercent = "referral_percent"
	CardIssueFee    = "card_issue_fee"
)
