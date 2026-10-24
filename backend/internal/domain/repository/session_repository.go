package repository

import (
	"context"

	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	logger "backend/internal/log"
	"backend/internal/tx"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (repository *SessionRepository) CreateSession(ctx context.Context, session entities.Sessions) (err error) {
	logger.Info("SessionRepository: CreateSession")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	if err = db.Create(&session).Error; err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
	}

	return
}
