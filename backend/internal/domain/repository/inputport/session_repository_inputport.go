//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock
package inputport

import (
	"context"

	"backend/internal/domain/entities"
)

type SessionRepositoryInputPort interface {
	CreateSession(ctx context.Context, session entities.Sessions) (err error)
	GetSessionByTokenHash(ctx context.Context, tokenHash string) (session *entities.Sessions, err error)
	DeleteSessionById(ctx context.Context, sessionId string) (err error)
}
