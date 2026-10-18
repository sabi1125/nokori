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

func (repository *AuthRepository) GetVerifiationCodeByEmail(ctx context.Context, userId string) (verificationCodes entities.VerificationCodes, err error) {
	logger.Info("AuthRepository: GetVerification")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	if err = db.Where("user_id = ?", userId).Take(&verificationCodes).Error; err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}
	return
}

func (repository *AuthRepository) UpdateVerificationCodeAttempts(ctx context.Context, verificationCodeId string, attempts int) (err error) {
	logger.Info("AuthRepository: UpdateVerificationCodeAttempts")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	if err = db.Model(&entities.VerificationCodes{}).Where("verification_code_id = ?", verificationCodeId).Update("attempts", attempts).Error; err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}
	return
}

func (repository *AuthRepository) UpdateVerificationCodeByVerificationCodeId(ctx context.Context, newVerificationCode entities.VerificationCodes) (err error) {
	logger.Info("AuthRepository: UpdateVerificationCodeByVerificationCodeId")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	if err = db.Model(&entities.VerificationCodes{}).Where("verification_code_id = ?", newVerificationCode.VerificationCodeId).
		Select("code", "resend_email_id", "expires_at", "attempts").
		Updates(newVerificationCode).Error; err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}

	return
}
