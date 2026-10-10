package interactor

import (
	"context"
	"time"

	api_inputport "backend/internal/domain/api_repository/inputport"
	"backend/internal/domain/apperror"
	"backend/internal/domain/entities"
	"backend/internal/domain/repository/inputport"
	logger "backend/internal/log"
	"backend/internal/tx"
	"backend/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthInteractor struct {
	authRepository                 inputport.AuthRepositoryInputPort
	userRepository                 inputport.UserRepositoryInputPort
	sendVerificationMailRepository api_inputport.SendVerificationMailInputPort
	uuidGenerator                  util.UUIDGenerator
	timeProvider                   util.TimeProvider
	txManager                      tx.Manager
}

func NewAuthInteractor(
	authRepository inputport.AuthRepositoryInputPort,
	userRepository inputport.UserRepositoryInputPort,
	sendVerificationMailRepository api_inputport.SendVerificationMailInputPort,
	uuidGenerator util.UUIDGenerator,
	timeProvider util.TimeProvider,
	txManager tx.Manager,
) *AuthInteractor {
	return &AuthInteractor{
		authRepository:                 authRepository,
		userRepository:                 userRepository,
		sendVerificationMailRepository: sendVerificationMailRepository,
		uuidGenerator:                  uuidGenerator,
		timeProvider:                   timeProvider,
		txManager:                      txManager,
	}
}

const verificationCodeLifetime = time.Minute

func (interactor *AuthInteractor) SignUp(ctx context.Context, signupParams entities.SignUp) (err error) {
	logger.Info("AuthInteractor: SignUp")

	user, err := interactor.userRepository.FindUserByEmail(ctx, signupParams.Email)
	if err != nil {
		return
	}

	if user.Email != "" && user.Verified {
		return new(apperror.UserAlreadyExists)
	}

	if user.Email != "" && !user.Verified {
		return new(apperror.EmailNotVerified)
	}

	newUserId, err := interactor.uuidGenerator.NewV7()
	if err != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupParams.Password), bcrypt.DefaultCost)
	if err != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	newUser := entities.Users{
		UserID:       newUserId,
		FirstName:    signupParams.FirstName,
		LastName:     signupParams.LastName,
		Email:        signupParams.Email,
		PasswordHash: string(hashedPassword),
		Verified:     false,
	}

	newVerificationCodeId, err := interactor.uuidGenerator.NewV7()
	if err != nil {
		return apperror.Wrap(apperror.InternalError, err)
	}

	resendId, verificationCode, err := interactor.sendVerificationMailRepository.SendEmailWithVerificationCode(newUser.Email)
	if err != nil {
		return
	}

	expiryTime := interactor.timeProvider.Now().Add(verificationCodeLifetime)

	verificationCodeParam := entities.VerificationCodes{
		VerificationCodeId: newVerificationCodeId,
		UserId:             newUserId,
		Code:               verificationCode,
		ResendEmailId:      resendId,
		ExpiresAt:          expiryTime,
	}

	txErr := interactor.txManager.WithinTransaction(ctx, func(ctx context.Context) error {
		if err := interactor.userRepository.CreateUser(ctx, newUser); err != nil {
			return err
		}

		if err := interactor.authRepository.InsertVerificationCode(ctx, verificationCodeParam); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		err = txErr
		return
	}

	return
}
