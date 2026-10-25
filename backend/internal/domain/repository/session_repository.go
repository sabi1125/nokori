package repository

import (
	"context"
	"errors"

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
		return
	}

	return
}

func (repository *SessionRepository) GetSessionByTokenHash(ctx context.Context, tokenHash string) (session *entities.Sessions, err error) {
	logger.Info("SessionRepository: GetSessionByTokenHash")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	var found entities.Sessions
	if err = db.Where("token_hash = ?", tokenHash).Take(&found).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = nil
			return
		}
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}

	session = &found

	return
}

func (repository *SessionRepository) DeleteSessionById(ctx context.Context, sessionId string) (err error) {
	logger.Info("SessionRepository: DeleteSessionById")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}
	res := db.Where("session_id = ?", sessionId).Delete(&entities.Sessions{})
	if res.Error != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	return
}
