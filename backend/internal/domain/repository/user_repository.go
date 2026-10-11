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

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) CreateUser(ctx context.Context, user entities.Users) (err error) {
	logger.Info("UserRepository: CreateUser")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	if err = db.Create(&user).Error; err != nil {
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}
	return
}

func (repository *UserRepository) FindUserByEmail(ctx context.Context, email string) (users entities.Users, err error) {
	logger.Info("UserRepository: FindUserByEmail")
	db := tx.ExtractTx(ctx)
	if db == nil {
		db = repository.db
	}

	if err = db.Where("email= ?", email).Take(&users).Error; err != nil {
		// if email address is not found return with no error
		// email not existing simply just means that email has never been registered
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("No such email address exists.")
			err = nil
			return
		}
		err = apperror.Wrap(apperror.InternalError, err)
		return
	}
	return
}
