package entities

import "time"

type VerificationCodes struct {
	VerificationCodeId string    `json:"verification_code_id" gorm:"column:verification_code_id"`
	UserId             string    `json:"user_id" grom:"column:user_id"`
	Code               string    `json:"code" gorm:"column:code"`
	ResendEmailId      string    `json:"resend_email_id" gorm:"column:resend_email_id"`
	ExpiresAt          time.Time `json:"expires_at" gorm:"column:expires_at"`
}
