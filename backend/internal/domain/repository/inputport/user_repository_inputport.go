package inputport

import (
	"context"

	"backend/internal/domain/entities"
)

type UserRepositoryInputPort interface {
	CreateUser(ctx context.Context, user entities.Users) (err error)
	FindUserByEmail(ctx context.Context, email string) (user entities.Users, err error)
}
