package repository

import (
	"context"

	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	logger "backend/internal/log"
	"backend/internal/tx"

	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repository *AuthRepository) InsertVerificationCode(ctx context.Context, verificationCode entities.VerificationCodes) (err error) {
	logger.Info("AuthRepository: InsertVerificationCode")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	if err = db.Create(&verificationCode).Error; err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}
	return
}
