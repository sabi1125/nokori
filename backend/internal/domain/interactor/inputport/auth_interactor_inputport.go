//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package inputport

import (
	"context"

	"backend/internal/domain/entities"
)

type AuthInteractorInputPort interface {
	SignUp(ctx context.Context, signupParams entities.SignUp) (err error)
}
