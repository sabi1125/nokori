//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package inputport

import (
	"context"

	"backend/internal/domain/entities"
)

type AuthRepositoryInputPort interface {
	InsertVerificationCode(ctx context.Context, verificationCode entities.VerificationCodes) (err error)
	GetVerifiationCodeByEmail(ctx context.Context, userId string) (verificationCodes entities.VerificationCodes, err error)
	UpdateVerificationCodeAttempts(ctx context.Context, verificationCodeId string, attempts int) (err error)
	UpdateVerificationCodeByVerificationCodeId(ctx context.Context, verificationCode entities.VerificationCodes) (err error)
}
