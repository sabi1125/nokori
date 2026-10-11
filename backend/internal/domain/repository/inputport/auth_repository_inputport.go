//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package inputport

import (
	"context"

	"backend/internal/domain/entities"
)

type AuthRepositoryInputPort interface {
	InsertVerificationCode(ctx context.Context, verificationCode entities.VerificationCodes) (err error)
}
