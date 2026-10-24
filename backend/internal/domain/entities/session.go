package entities

import "time"

type Sessions struct {
	SessionId          string     `json:"session_id" gorm:"column:session_id"`
	UserId             string     `json:"user_id" gorm:"column:user_id"`
	TokenHash          string     `json:"token_hash" gorm:"column:token_hash"`
	PreviousTokenHash  *string    `json:"previous_token_hash" gorm:"column:previous_token_hash"`
	PreviousValidUntil *time.Time `json:"previous_valid_until" gorm:"column:previous_valid_until"`
	ExpiresAt          time.Time  `json:"expires_at" gorm:"column:expires_at"`
}
